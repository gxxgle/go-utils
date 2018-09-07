package math

import (
	"testing"
)

func TestLog(t *testing.T) {
	products := []float64{
		10,
		100,
		0.1,
		0.1,
		0.001,
		0.0001,
	}

	bases := []float64{
		10,
		10,
		0.1,
		1,
		0.998,
		0.998,
	}

	exes := []float64{
		1,
		2,
		1,
		0,
		3451,
		4601,
	}

	for i, product := range products {
		out := Log(product, bases[i])
		if out != exes[i] {
			t.Errorf("Log(%v, %v) != %v, but %v", product, bases[i], exes[i], out)
		}
	}
}
