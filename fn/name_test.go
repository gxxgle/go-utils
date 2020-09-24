package fn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type s struct{}

func f1(int)        {}
func f2(...int)     {}
func F3(int)        {}
func F4(...int)     {}
func (s) f1(int)    {}
func (s) f2(...int) {}
func (s) F3(int)    {}
func (s) F4(...int) {}

func TestName(t *testing.T) {
	type check struct {
		input    interface{}
		expected string
		err      error
	}

	checks := []check{
		{input: f1, expected: "f1", err: nil},
		{input: f2, expected: "f2", err: nil},
		{input: F3, expected: "F3", err: nil},
		{input: F4, expected: "F4", err: nil},
		{input: s{}.f1, expected: "f1", err: nil},
		{input: s{}.f2, expected: "f2", err: nil},
		{input: s{}.F3, expected: "F3", err: nil},
		{input: s{}.F4, expected: "F4", err: nil},
		{input: 1, expected: "", err: fmt.Errorf("int is not func type")},
		{input: true, expected: "", err: fmt.Errorf("bool is not func type")},
	}

	for _, c := range checks {
		output, err := Name(c.input)
		require.Equal(t, c.expected, output)
		require.Equal(t, c.err, err)
	}
}
