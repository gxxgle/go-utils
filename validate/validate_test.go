package validate

import (
	"testing"
)

func TestTagDefault(t *testing.T) {
	type myStruct struct {
		String1 string  `validate:"default=127.0.0.1"`
		String2 string  `validate:"default=127.0.0.1"`
		Int     int     `validate:"default=8080"`
		Float   float64 `validate:"default=3.14"`
	}

	my := &myStruct{String1: "string"}
	if err := V.Struct(my); err != nil {
		t.Fatal(err)
	}

	if my.String1 != "string" {
		t.Fatal("TestTagDefault String1")
	}

	if my.String2 != "127.0.0.1" {
		t.Fatal("TestTagDefault String2")
	}

	if my.Int != 8080 {
		t.Fatal("TestTagDefault Int")
	}

	if my.Float != 3.14 {
		t.Fatal("TestTagDefault Float")
	}
}

func TestTagEmpty(t *testing.T) {
	type myStruct struct {
		Bind string
	}

	my := &myStruct{}
	if err := V.Struct(my); err != nil {
		t.Fatal("TestTagEmpty failed")
	}
}

func TestExample(t *testing.T) {
	type myStruct struct {
		Bind     string `validate:"default=0.0.0.0:0"`
		IP       string `validate:"default=127.0.0.1,ip"`
		PoolSize int    `validate:"default=50,gte=1,lte=200"`
		Email    string `validate:"email"`
		Username string `validate:"min=6,max=12"`
		Gender   string `validate:"oneof=M FM"`
	}

	my1 := &myStruct{}
	if err := V.Struct(my1); err == nil {
		t.Fatal("TestExample my1 failed")
	}

	my2 := &myStruct{
		Email:    "admin@abc.xyz",
		Username: "username",
		Gender:   "FM",
	}
	if err := V.Struct(my2); err != nil {
		t.Fatal("TestExample my2 err:", err)
	}
}

func TestDeep(t *testing.T) {
	type Config struct {
		Bind     string `validate:"default=0.0.0.0:0"`
		IP       string `validate:"default=127.0.0.1,ip"`
		PoolSize int    `validate:"default=50,gte=1,lte=200"`
		Email    string `validate:"email"`
		Username string `validate:"min=6,max=12"`
		Gender   string `validate:"oneof=M FM"`
	}

	type myStruct struct {
		Redis *Config `json:"redis" yaml:"redis" validate:"required"`
	}

	my1 := &myStruct{}
	if err := V.Struct(my1); err == nil {
		t.Fatal("TestDeep my1 failed")
	}

	my2 := &myStruct{
		Redis: &Config{
			Email:    "admin@abc.xyz",
			Username: "username",
			Gender:   "FM",
		},
	}
	if err := V.Struct(my2); err != nil {
		t.Fatal("TestDeep my2 err:", err)
	}
}
