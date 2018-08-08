package schedule

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron"
)

// global variable
var (
	Cron *cron.Cron
	wg   sync.WaitGroup
)

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}

	Cron = cron.NewWithLocation(loc)
	Cron.Start()
}

// AddFunc [docs](https://godoc.org/github.com/robfig/cron)
func AddFunc(spec string, fn func()) error {
	return Cron.AddFunc(spec, func() {
		wg.Add(1)
		defer wg.Done()

		fn()
	})
}

func Close() {
	Cron.Stop()
	wg.Wait()
}
