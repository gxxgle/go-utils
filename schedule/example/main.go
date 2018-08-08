package main

import (
	"log"
	"time"

	"github.com/gxxgle/go-utils/schedule"
)

func do() {
	log.Println("task begin")
	time.Sleep(time.Second * 10)
	log.Println("task end")
}

func main() {
	err := schedule.AddFunc("@every 1s", do)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 5)
	schedule.Close()
	log.Println("stopped")
}
