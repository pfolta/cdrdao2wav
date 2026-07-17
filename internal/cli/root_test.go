package cli

import (
	"bytes"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestNewRootCommand(t *testing.T) {

	// Run the root command with the built-in "help" subcommand.
	helpOut := new(bytes.Buffer)
	cmd := NewRootCommand()
	cmd.SetOut(helpOut)
	cmd.SetArgs([]string{"help"})
	err := cmd.Execute()

	// The command should print the help text and return successfully.
	assert.NilError(t, err)
	assert.Assert(t, is.Contains(helpOut.String(), "Usage:"))

	// Run the root command without arguments.
	rootOut := new(bytes.Buffer)
	cmd = NewRootCommand()
	cmd.SetOut(rootOut)
	cmd.SetArgs([]string{})
	err = cmd.Execute()

	// The command should print the same help text but return an ErrNoCommand.
	assert.ErrorIs(t, err, ErrNoCommand)
	assert.Equal(t, rootOut.String(), helpOut.String())
}
