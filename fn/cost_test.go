package fn

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCost(t *testing.T) {
	fn := func() { time.Sleep(time.Second) }
	require.Equal(t, time.Second, Cost(fn).Truncate(time.Millisecond*10))
}
