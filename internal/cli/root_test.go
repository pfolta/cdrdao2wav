// Copyright (c) 2026 Peter Folta
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
