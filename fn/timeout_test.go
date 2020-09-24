package fn

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stringErr string

func (e stringErr) Error() string {
	return string(e)
}

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
			input:    genFunc(time.Millisecond, stringErr("func err")),
			timeout:  time.Millisecond * 2,
			expected: stringErr("func err"),
		},
		{
			input:    genFunc(time.Millisecond*2, stringErr("func err")),
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
