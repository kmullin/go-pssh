package ssh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSHOption(t *testing.T) {
	testCases := []struct {
		key      string
		value    any
		expected string
	}{
		{"foobar", 15, "-o foobar=15"},
		{"something", true, "-o something=yes"},
		{"baz", false, "-o baz=no"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			assert.Equal(t, tc.expected, sshOption(tc.key, tc.value))
		})
	}
}
