package nhttp_test

import (
	"github.com/smartwalle/nhttp"
	"net/url"
	"testing"
)

type Human struct {
	Name     string `form:"name"`
	Age      int    `form:"age"`
	Birthday string `form:"birthday"`
}

func TestMapper_BindHuman(t *testing.T) {
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
	}

	var m = nhttp.NewMapper("form")

	for _, test := range tests {
		var dst Human

		if err := m.Bind(test.form, &dst); err != nil {
			t.Fatal(err)
		}

		if dst.Name != test.human.Name {
			t.Fatalf("Name - 期望: %s, 实际: %s", test.human.Name, dst.Name)
		}

		if dst.Age != test.human.Age {
			t.Fatalf("Age - 期望: %d, 实际: %d", test.human.Age, dst.Age)
		}

		if dst.Birthday != test.human.Birthday {
			t.Fatalf("Age - 期望: %s, 实际: %s", test.human.Birthday, dst.Birthday)
		}
	}
}

type Student struct {
	Human
	Number int64  `form:"number"`
	Class  string `form:"class"`
}

var studentForm = url.Values{"name": {"Yangfeng"}, "age": {"10"}, "number": {"3414416614257328130"}, "birthday": {"2016-06-12"}, "class": {"class"}}

func TestMapper_Bind(t *testing.T) {
	var s *Student
	var m = nhttp.NewMapper("form")
	m.Bind(studentForm, &s)
	t.Log(s)
}

func BenchmarkMapperBind(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var s Student
		m.Bind(studentForm, &s)
	}
	b.StopTimer()
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
		m.Bind(userForm, &u)
	}
}

func BenchmarkMapperBind_Pointer(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var u *User
		m.Bind(userForm, &u)
	}
}

func BenchmarkMapperBind_StructParallel(b *testing.B) {
	var m = nhttp.NewMapper("form")
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var u User
			m.Bind(userForm, &u)
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
			m.Bind(userForm, &u)
		}
	})
}
