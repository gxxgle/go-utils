package fn

import (
	"time"
)

func Cost(fn func()) (cost time.Duration) {
	startAt := time.Now()

	defer func() {
		cost = time.Since(startAt)
	}()

	fn()
	return cost
}
