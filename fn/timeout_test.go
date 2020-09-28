package fn

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	genFunc := func(d time.Duration, err error) func() error {
		return func() error {
			time.Sleep(d)
			return err
		}
	}

	type check struct {
		input    func() error
		timeout  time.Duration
		expected error
	}

	checks := []check{
		{
			input:    genFunc(time.Millisecond, nil),
			timeout:  time.Millisecond * 2,
			expected: nil,
		},
		{
			input:    genFunc(time.Millisecond, errors.New("func err")),
			timeout:  time.Millisecond * 2,
			expected: errors.New("func err"),
		},
		{
			input:    genFunc(time.Millisecond*2, errors.New("func err")),
			timeout:  time.Millisecond,
			expected: context.DeadlineExceeded,
		},
	}

	for i, c := range checks {
		err := Timeout(context.Background(), c.input, c.timeout)
		require.Equal(t, c.expected, err, "#%d check", i)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()
	require.Equal(t, context.Canceled, Timeout(ctx, genFunc(time.Minute, nil), time.Hour))
}
