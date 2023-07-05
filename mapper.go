package nhttp

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	kNoTag   = "-"
	kTag     = "form"
	kDefault = "default"
)

type DecodeFunc func(name, tag string, values []string) (interface{}, error)

type fieldDescriptor struct {
	Index   []int
	Tag     string
	Default []string
	Decoder DecodeFunc
}

type structDescriptor struct {
	Name   string
	Fields []fieldDescriptor
}

type Mapper struct {
	tag      string
	structs  atomic.Value // map[reflect.Type]structDescriptor
	mu       sync.Mutex
	decoders map[reflect.Type]DecodeFunc
}

var mapper = NewMapper(kTag)

func Bind(src map[string][]string, dst interface{}) error {
	return mapper.Bind(src, dst)
}

func NewMapper(tag string) *Mapper {
	var m = &Mapper{}
	m.tag = tag
	m.structs.Store(make(map[reflect.Type]structDescriptor))
	m.decoders = make(map[reflect.Type]DecodeFunc)
	return m
}

func (this *Mapper) UseDecoder(dstType reflect.Type, fn DecodeFunc) {
	if dstType == nil || fn == nil {
		return
	}
	this.decoders[dstType] = fn
}

func (this *Mapper) Bind(src map[string][]string, dst interface{}) error {
	var dstValue = reflect.ValueOf(dst)
	var dstType = dstValue.Type()

	if dstValue.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer to Bind")
	}

	if dstValue.IsNil() {
		return errors.New("nil pointer passed to Bind")
	}

	for {
		if dstValue.Kind() == reflect.Ptr && dstValue.IsNil() {
			dstValue.Set(reflect.New(dstType.Elem()))
		}

		if dstValue.Kind() == reflect.Ptr {
			dstValue = dstValue.Elem()
			dstType = dstType.Elem()
			continue
		}
		break
	}

	var dStruct, ok = this.getStructDescriptor(dstType)
	if !ok {
		dStruct = this.parseStructDescriptor(dstType)
	}

	for _, field := range dStruct.Fields {
		var values, exists = src[field.Tag]
		if !exists {
			if len(field.Default) == 0 {
				continue
			}
			values = field.Default
		}

		var fieldValue = fieldByIndex(dstValue, field.Index)

		if field.Decoder != nil {
			var nValue, err = field.Decoder(dStruct.Name, field.Tag, values)
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(nValue))
			continue
		}

		if err := mapValues(fieldValue, values); err != nil {
			return err
		}
	}
	return nil
}

func fieldByIndex(parent reflect.Value, index []int) reflect.Value {
	if len(index) == 1 {
		return parent.Field(index[0])
	}
	for i, x := range index {
		if i > 0 {
			if parent.Kind() == reflect.Pointer && parent.Type().Elem().Kind() == reflect.Struct {
				if parent.IsNil() {
					parent.Set(reflect.New(parent.Type().Elem()))
				}
				parent = parent.Elem()
			}
		}
		parent = parent.Field(x)
	}
	return parent
}

func (this *Mapper) getStructDescriptor(key reflect.Type) (structDescriptor, bool) {
	var value, ok = this.structs.Load().(map[reflect.Type]structDescriptor)[key]
	return value, ok
}

func (this *Mapper) setStructDescriptor(key reflect.Type, value structDescriptor) {
	var structs = this.structs.Load().(map[reflect.Type]structDescriptor)
	structs[key] = value
	this.structs.Store(structs)
}

type structQueueElement struct {
	Type  reflect.Type
	Index []int
}

func (this *Mapper) parseStructDescriptor(dstType reflect.Type) structDescriptor {
	this.mu.Lock()

	var dStruct, ok = this.getStructDescriptor(dstType)
	if ok {
		this.mu.Unlock()
		return dStruct
	}

	var queue = make([]structQueueElement, 0, 3)
	queue = append(queue, structQueueElement{
		Type:  dstType,
		Index: nil,
	})

	var dFields = make(map[string]fieldDescriptor)

	for len(queue) > 0 {
		var current = queue[0]
		queue = queue[1:]

		var numField = current.Type.NumField()

		for i := 0; i < numField; i++ {
			var fieldStruct = current.Type.Field(i)

			var tag = fieldStruct.Tag.Get(this.tag)
			if tag == kNoTag {
				continue
			}

			var opts string
			tag, opts = head(tag, ",")

			if tag == "" {
				tag = fieldStruct.Name

				if fieldStruct.Type.Kind() == reflect.Ptr {
					queue = append(queue, structQueueElement{
						Type:  fieldStruct.Type.Elem(),
						Index: append(current.Index, i),
					})
					continue
				}

				if fieldStruct.Type.Kind() == reflect.Struct {
					queue = append(queue, structQueueElement{
						Type:  fieldStruct.Type,
						Index: append(current.Index, i),
					})
					continue
				}
			}

			if _, exists := dFields[tag]; exists {
				continue
			}

			var dField = fieldDescriptor{}
			dField.Index = append(current.Index, i)
			dField.Tag = tag
			dField.Decoder = this.decoders[fieldStruct.Type]

			var opt string
			for len(opts) > 0 {
				opt, opts = head(opts, ",")

				key, value := head(opt, "=")
				switch key {
				case kDefault:
					dField.Default = strings.Split(value, ";")
				}
			}

			dFields[tag] = dField
		}
	}

	dStruct.Name = dstType.String()
	dStruct.Fields = make([]fieldDescriptor, 0, len(dFields))
	for _, field := range dFields {
		dStruct.Fields = append(dStruct.Fields, field)
	}

	this.setStructDescriptor(dstType, dStruct)
	this.mu.Unlock()

	return dStruct
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

func mapValues(field reflect.Value, values []string) error {
	if field.Kind() == reflect.Slice {
		var valueLen = len(values)
		var s = reflect.MakeSlice(field.Type(), valueLen, valueLen)
		for i := 0; i < valueLen; i++ {
			if err := mapValue(s.Index(i), values[i]); err != nil {
				return err
			}
		}
		field.Set(s)
		return nil
	}
	return mapValue(field, values[0])
}

func mapValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Interface:
		field.Set(reflect.ValueOf(value))
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		return mapInt(field, value, 0)
	case reflect.Int8:
		return mapInt(field, value, 8)
	case reflect.Int16:
		return mapInt(field, value, 16)
	case reflect.Int32:
		return mapInt(field, value, 32)
	case reflect.Int64:
		return mapInt(field, value, 64)
	case reflect.Uint:
		return mapUint(field, value, 0)
	case reflect.Uint8:
		return mapUint(field, value, 8)
	case reflect.Uint16:
		return mapUint(field, value, 16)
	case reflect.Uint32:
		return mapUint(field, value, 32)
	case reflect.Uint64:
		return mapUint(field, value, 64)
	case reflect.Uintptr:
		return mapUint(field, value, 64)
	case reflect.Float32:
		return mapFloat(field, value, 32)
	case reflect.Float64:
		return mapFloat(field, value, 64)
	case reflect.Bool:
		return mapBool(field, value)
	default:
		return errors.New("cannot unmarshal into " + field.Type().String())
	}
	return nil
}

func mapInt(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		field.SetInt(0)
		return nil
	}
	intValue, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return err
	}
	field.SetInt(intValue)
	return nil
}

func mapUint(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		field.SetUint(0)
		return nil
	}
	intValue, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return err
	}
	field.SetUint(intValue)
	return nil
}

func mapFloat(field reflect.Value, value string, bitSize int) error {
	if value == "" {
		field.SetFloat(0)
		return nil
	}
	intValue, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return err
	}
	field.SetFloat(intValue)
	return nil
}

func mapBool(field reflect.Value, value string) error {
	if value == "" {
		field.SetBool(false)
		return nil
	}
	booleValue, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(booleValue)
	return nil
}
