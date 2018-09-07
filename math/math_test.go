package math

import (
	"testing"
)

func TestFloatToInt(t *testing.T) {
	ins := []float64{
		15.666600002,
		15.666599998,
	}

	precisions := []int32{
		4,
		4,
	}

	exes := []int{
		156666,
		156666,
	}

	for i, in := range ins {
		out := FloatToInt(in, precisions[i])
		if out != exes[i] {
			t.Errorf("FloatToInt(%v, %v) != %v, but %v", in, precisions[i], exes[i], out)
		}
	}
}

func TestFloatFromInt(t *testing.T) {
	ins := []int{
		156666,
		156666,
	}

	precisions := []int32{
		4,
		4,
	}

	exes := []float64{
		15.6666,
		15.6666,
	}

	for i, in := range ins {
		out := FloatFromInt(in, precisions[i])
		if out != exes[i] {
			t.Errorf("FloatFromInt(%v, %v) != %v, but %v", in, precisions[i], exes[i], out)
		}
	}
}
