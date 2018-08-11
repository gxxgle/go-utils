package schedule

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron"
)

// global variable
var (
	Cron    *cron.Cron
	stopped bool
	wg      sync.WaitGroup
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}

	Cron = cron.NewWithLocation(loc)
	Cron.Start()
}

// AddCronFunc [docs](https://godoc.org/github.com/robfig/cron)
func AddCronFunc(spec string, fn func()) error {
	return Cron.AddFunc(spec, func() {
		wg.Add(1)
		defer wg.Done()

		fn()
	})
}

// AddLoopFunc can loop run task.
func AddLoopFunc(sleep time.Duration, fn func()) {
	go func() {
		for !stopped {
			wg.Add(1)
			fn()
			wg.Done()

			time.Sleep(sleep)
		}
	}()
}

func Close() {
	stopped = true
	Cron.Stop()
	wg.Wait()
}
