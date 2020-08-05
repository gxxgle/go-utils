package math

import (
	"math"
)

// Log args: [product, base(e), step(e)]
func Log(args ...float64) float64 {
	var (
		multi   = float64(1)
		product = float64(1)
		base    = float64(math.E)
		step    = float64(1)
	)

	if len(args) == 0 {
		return 0
	}

	if len(args) > 0 {
		product = args[0]
	}

	if len(args) > 1 {
		base = args[1]
	}

	if base < 0 || base == 1 {
		return 0
	}

	if base < 1 {
		multi = float64(-1)
	}

	if len(args) > 2 {
		step = args[2]
	}

	if step <= 0 {
		return 0
	}

	cout := 1
	i := float64(0.0)
	for {
		if cout >= 100*10000 {
			return 0
		}

		if math.Pow(base, i)*multi >= product*multi {
			return i
		}

		cout++
		i += step
	}
}
