package fn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParallel(t *testing.T) {
	a1, a2, a3 := 0, 0, 0
	f1 := func() {
		time.Sleep(time.Second * 1)
		a1 = 1
	}
	f2 := func() {
		time.Sleep(time.Second * 2)
		a2 = 2
	}
	f3 := func() {
		time.Sleep(time.Second * 3)
		a3 = 3
	}

	cost := Cost(func() { Parallel(f1, f2, f3) })
	require.Equal(t, 1, a1)
	require.Equal(t, 2, a2)
	require.Equal(t, 3, a3)
	require.Equal(t, time.Second*3, cost.Truncate(time.Millisecond*10))
}
