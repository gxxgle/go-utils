package fn

import (
	"context"
	"time"
)

func Timeout(ctx context.Context, fn func() error, timeout time.Duration) error {
	errChan := make(chan error)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		errChan <- fn()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
