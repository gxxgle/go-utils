package main

import (
	sn "sync"
	"time"

	"github.com/gxxgle/go-utils/cache"
	"github.com/gxxgle/go-utils/log"
	"github.com/gxxgle/go-utils/sync"
)

var (
	wg sn.WaitGroup
)

func init() {
	log.LogIfError(sync.InitRedis(&cache.RedisConfig{
		URL:      "devhost:6379",
		Password: "KgqvdOdYV5",
		Retries:  10,
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
	log.L.Info("lock getting index:", i, ", cost:", time.Since(begin).String())
	mu.Lock()
	log.L.Info("lock got index:", i, ", cost:", time.Since(begin).String())
	defer wg.Done()
	defer mu.Unlock()
	time.Sleep(time.Second * 5)
	log.L.Info("lock un index:", i, ", cost:", time.Since(begin).String())
}
