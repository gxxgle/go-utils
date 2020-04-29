package main

import (
	"github.com/gxxgle/go-utils/log"
)

type person struct {
	Name string
	Age  int
}

func (p *person) String() string {
	return "xxx"
}

func main() {
	log.File()
	log.Debug()
	log.L.WithFields(log.F{
		"c":       "c",
		"z":       "z",
		"a":       "a",
		"b":       "b",
		"t":       "t",
		"persion": &person{},
	}).Debug("test")
}
