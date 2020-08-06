package main

import (
	"time"

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
	log.TextFormat()
	// log.File()
	log.Debug()
	fields := log.F{
		"c":       "c",
		"z":       "z",
		"a":       "a",
		"b":       "b",
		"t":       "t",
		"persion": &person{},
	}
	log.L.WithFields(fields).Debug("test debug")
	time.Sleep(time.Millisecond * 347)
	log.L.WithFields(fields).Info("test info")
	time.Sleep(time.Millisecond * 347)
	log.L.WithFields(fields).Warn("test warn")
	time.Sleep(time.Millisecond * 347)
	log.L.WithFields(fields).Error("test error")
}
