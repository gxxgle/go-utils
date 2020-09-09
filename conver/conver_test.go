package conver

import (
	"errors"
	"fmt"
	"testing"
)

type intStr int

func (i intStr) String() string {
	return fmt.Sprintf("intStr(%d)", i)
}

type myStruct struct {
	Name string
}

func TestString(t *testing.T) {
	type check struct {
		val  interface{}
		want string
	}

	checks := []check{
		{val: "str", want: "str"},
		{val: []byte("bytes"), want: "bytes"},
		{val: intStr(1), want: "intStr(1)"},
		{val: errors.New("err"), want: "err"},
		{val: 123, want: "123"},
		{val: 3.14, want: "3.14"},
		{val: myStruct{Name: "name"}, want: `conver.myStruct{Name:"name"}`},
		{val: []int64{1, 2, 3}, want: "[]int64{1, 2, 3}"},
	}

	for _, c := range checks {
		str, err := String(c.val)
		if err != nil {
			t.Fatal("TestString err:", err)
		}
		if str != c.want {
			t.Fatalf("TestString(%v) want: %s, but: %s", c.val, c.want, str)
		}
	}
}

func TestBool(t *testing.T) {
	type check struct {
		val  interface{}
		want bool
	}

	checks := []check{
		{val: true, want: true},
		{val: false, want: false},
		{val: 1, want: true},
		{val: 0, want: false},
		{val: 0.0, want: false},
		{val: "1", want: true},
		{val: "t", want: true},
		{val: "T", want: true},
		{val: "true", want: true},
		{val: "TRUE", want: true},
		{val: "True", want: true},
		{val: "ok", want: true},
		{val: "OK", want: true},
		{val: "yes", want: true},
		{val: "YES", want: true},
		{val: "0", want: false},
		{val: "f", want: false},
		{val: "F", want: false},
		{val: "false", want: false},
		{val: "FALSE", want: false},
		{val: "False", want: false},
		{val: "", want: false},
	}

	for i, c := range checks {
		bl, err := Bool(c.val)
		if err != nil {
			t.Fatal("TestBool err:", err)
		}
		if bl != c.want {
			t.Fatalf("%d TestBool(%v) want: %v, but: %v", i, c.val, c.want, bl)
		}
	}
}

func TestFloat64(t *testing.T) {
	type check struct {
		val  interface{}
		want float64
	}

	checks := []check{
		{val: true, want: 1},
		{val: false, want: 0},
		{val: 1, want: 1},
		{val: -1, want: -1},
		{val: "1", want: 1},
		{val: "-1", want: -1},
		{val: "3.1415", want: 3.1415},
		{val: "-3.1415", want: -3.1415},
	}

	for i, c := range checks {
		fl, err := Float64(c.val)
		if err != nil {
			t.Fatal("TestFloat64 err:", err)
		}
		if fl != c.want {
			t.Fatalf("%d TestFloat64(%v) want: %v, but: %v", i, c.val, c.want, fl)
		}
	}
}

func TestInt64(t *testing.T) {
	type check struct {
		val  interface{}
		want int64
	}

	checks := []check{
		{val: true, want: 1},
		{val: false, want: 0},
		{val: 1, want: 1},
		{val: -1, want: -1},
		{val: "1", want: 1},
		{val: "-1", want: -1},
		{val: "3.1415", want: 3},
		{val: "-3.1415", want: -3},
		{val: "3.91415", want: 3},
		{val: "-3.91415", want: -3},
	}

	for i, c := range checks {
		it, err := Int64(c.val)
		if err != nil {
			t.Fatal("TestInt64 err:", err)
		}
		if it != c.want {
			t.Fatalf("%d TestInt64(%v) want: %v, but: %v", i, c.val, c.want, it)
		}
	}
}
