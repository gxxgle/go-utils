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
	// cron()
	loop()
}

func cron() {
	err := schedule.AddCronFunc("@every 1s", do)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 5)
	log.Println("stop begin")
	schedule.Close()
	log.Println("stop finish")
}

func loop() {
	schedule.AddLoopFunc(time.Second, do)
	time.Sleep(time.Second * 5)
	log.Println("stop begin")
	schedule.Close()
	log.Println("stop finish")
}
