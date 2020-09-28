package fn

import (
	"sync"
)

func Parallel(fns ...func()) {
	wg := sync.WaitGroup{}

	for _, fn := range fns {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(fn)
	}

	wg.Wait()
}
