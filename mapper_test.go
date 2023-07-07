package nhttp_test

import (
	"github.com/smartwalle/nhttp"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type Human struct {
	Name     string `form:"name"`
	Age      int    `form:"age,default=20"`
	Birthday string `form:"birthday"`
}

func TestMapper_Bind(t *testing.T) {
	var tests = []struct {
		form  url.Values
		human Human
	}{
		{
			form:  url.Values{"name": {"name1"}, "age": {"10"}, "birthday": {"2023-07-05"}},
			human: Human{Name: "name1", Age: 10, Birthday: "2023-07-05"},
		},
		{
			form:  url.Values{"name": {"name2"}, "age": {"11"}, "birthday": {"2023-07-06"}},
			human: Human{Name: "name2", Age: 11, Birthday: "2023-07-06"},
		},
		{
			form:  url.Values{"name": {"name3"}, "age": {"0"}, "birthday": {"2023-07-07"}},
			human: Human{Name: "name3", Age: 0, Birthday: "2023-07-07"},
		},
		{
			form:  url.Values{"name": {"name4"}, "age": {"-1"}, "birthday": {"2023-07-08"}},
			human: Human{Name: "name4", Age: -1, Birthday: "2023-07-08"},
		},
		{
			form:  url.Values{"name": {"name5"}, "age": {""}, "birthday": {"2023-07-09"}},
			human: Human{Name: "name5", Age: 0, Birthday: "2023-07-09"},
		},
		{
			form:  url.Values{"name": {"name6"}, "birthday": {"2023-07-10"}},
			human: Human{Name: "name6", Age: 20, Birthday: "2023-07-10"},
		},
	}

	for _, test := range tests {
		var dst Human

		if err := nhttp.Bind(test.form, &dst); err != nil {
			t.Fatal(err)
		}

		if dst.Name != test.human.Name {
			t.Fatalf("Name - 期望: %s, 实际: %s", test.human.Name, dst.Name)
		}

		if dst.Age != test.human.Age {
			t.Fatalf("Age - 期望: %d, 实际: %d", test.human.Age, dst.Age)
		}

		if dst.Birthday != test.human.Birthday {
			t.Fatalf("Birthday - 期望: %s, 实际: %s", test.human.Birthday, dst.Birthday)
		}
	}
}

type Student struct {
	Human
	Number int64  `form:"number"`
	Class  string `form:"class"`
}

func TestMapper_BindEmbed(t *testing.T) {
	var tests = []struct {
		form    url.Values
		student Student
	}{
		{
			form:    url.Values{"name": {"name1"}, "age": {"10"}, "birthday": {"2023-07-05"}, "number": {"1"}, "class": {"c1"}},
			student: Student{Human: Human{Name: "name1", Age: 10, Birthday: "2023-07-05"}, Number: 1, Class: "c1"},
		},
		{
			form:    url.Values{"name": {"name2"}, "age": {"11"}, "birthday": {"2023-07-06"}, "number": {"2"}, "class": {"c2"}},
			student: Student{Human: Human{Name: "name2", Age: 11, Birthday: "2023-07-06"}, Number: 2, Class: "c2"},
		},
	}

	for _, test := range tests {
		var dst Student

		if err := nhttp.Bind(test.form, &dst); err != nil {
			t.Fatal(err)
		}

		if dst.Name != test.student.Name {
			t.Fatalf("Name - 期望: %s, 实际: %s", test.student.Name, dst.Name)
		}

		if dst.Age != test.student.Age {
			t.Fatalf("Age - 期望: %d, 实际: %d", test.student.Age, dst.Age)
		}

		if dst.Birthday != test.student.Birthday {
			t.Fatalf("Birthday - 期望: %s, 实际: %s", test.student.Birthday, dst.Birthday)
		}

		if dst.Number != test.student.Number {
			t.Fatalf("Birthday - 期望: %d, 实际: %d", test.student.Number, dst.Number)
		}

		if dst.Class != test.student.Class {
			t.Fatalf("Birthday - 期望: %s, 实际: %s", test.student.Class, dst.Class)
		}
	}
}

var userForm = url.Values{"firstname": {"Feng"}, "lastname": {"Yang"}, "email": {"yangfeng@qq.com"}, "age": {"10"}}

type User struct {
	Firstname string `form:"firstname"`
	Lastname  string `form:"lastname"`
	Email     string `form:"email"`
	Age       int    `form:"age"`
}

func BenchmarkMapperBind_Struct(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var u User
		if err := m.Bind(userForm, &u); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapperBind_Pointer(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var u *User
		if err := m.Bind(userForm, &u); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapperBind_StructParallel(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var u User
			if err := m.Bind(userForm, &u); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMapperBind_PointerParallel(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var u *User
			if err := m.Bind(userForm, &u); err != nil {
				b.Fatal(err)
			}
		}
	})
}

var form = url.Values{"end_time": {"2022-02-02"}}

type DateTime struct {
	BeginTime *time.Time `form:"begin_time,default=2022-01-02"`
	EndTime   *time.Time `form:"end_time"`
}

func TestMapper_UseDecoder(t *testing.T) {
	var m = nhttp.NewMapper("form")
	m.UseDecoder(reflect.TypeOf(&time.Time{}), func(value string) (interface{}, error) {
		var nt, err = time.Parse("2006-01-02", value)
		if err != nil {
			return nil, err
		}
		return &nt, nil
	})

	var dst DateTime
	if err := m.Bind(form, &dst); err != nil {
		t.Fatal(err)
	}

	if dst.BeginTime == nil || dst.EndTime == nil {
		t.Fatal("转换时间失败")
	}

	if dst.BeginTime.Year() != 2022 || dst.BeginTime.Month() != time.January || dst.BeginTime.Day() != 2 {
		t.Fatal("转换时间失败")
	}

	if dst.EndTime.Year() != 2022 || dst.EndTime.Month() != time.February || dst.EndTime.Day() != 2 {
		t.Fatal("转换时间失败")
	}
}

func BenchmarkMapperUseDecoder(b *testing.B) {
	var m = nhttp.NewMapper("form")
	m.UseDecoder(reflect.TypeOf(&time.Time{}), func(value string) (interface{}, error) {
		var nt, err = time.Parse("2006-01-02", value)
		if err != nil {
			return nil, err
		}
		return &nt, nil
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dst DateTime
		if err := m.Bind(form, &dst); err != nil {
			b.Fatal(err)
		}
	}
}
