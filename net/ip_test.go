package net

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivateIP(t *testing.T) {
	type check struct {
		input    string
		expected bool
	}

	checks := []check{
		{input: "", expected: false},
		{input: "0.0.0.0", expected: false},
		{input: "127.0.0.1", expected: false},
		{input: "192.168.0.0", expected: true},
		{input: "192.168.255.255", expected: true},
		{input: "192.169.0.0", expected: false},
	}

	for _, c := range checks {
		assert.Equal(t, c.expected, IsPrivateIP(net.ParseIP(c.input)), c.input)
	}
}
