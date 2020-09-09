package main

import (
	sn "sync"
	"time"

	"github.com/gxxgle/go-utils/cache"
	ulog "github.com/gxxgle/go-utils/log"
	"github.com/gxxgle/go-utils/sync"

	"github.com/phuslu/log"
)

var (
	wg sn.WaitGroup
)

func init() {
	ulog.LogIfError(sync.InitRedis(&cache.RedisConfig{
		Addr:       "devhost:6379",
		Password:   "KgqvdOdYV5",
		MaxRetries: 10,
	}))
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go print(i + 1)
	}

	wg.Wait()
}

func print(i int) {
	begin := time.Now()
	mu := sync.NewMutex("LOCKER:TEST")
	log.Info().Int("index", i).Dur("cost", time.Since(begin)).Msg("lock getting")
	mu.Lock()
	log.Info().Int("index", i).Dur("cost", time.Since(begin)).Msg("lock got")
	defer wg.Done()
	defer mu.Unlock()
	time.Sleep(time.Second * 5)
	log.Info().Int("index", i).Dur("cost", time.Since(begin)).Msg("lock un")
}
