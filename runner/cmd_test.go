package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	r := New([]string{"sleep", "15"}, 1)

	cmd := r.newCmd("hostA")

	assert.Nil(t, cmd.cmd.Stdin)
	assert.Nil(t, cmd.cmd.Stdout)
	assert.Nil(t, cmd.cmd.Stderr)
	t.Logf("%#v", cmd.cmd)
}
