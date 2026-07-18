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

	"github.com/pfolta/cdrdao2audio"
)

var testMetadata = cdrdao2audio.Metadata{
	Name:      "cdrdao2audio",
	Version:   "v1.4.2-test",
	BuildDate: "2026-07-15T22:29:10Z",
	License:   "MIT",
	OS:        "darwin",
	Arch:      "arm64",
}

const (
	expectedShortVersion = "1.4.2-test\n"
	expectedLongVersion  = "cdrdao2audio version v1.4.2-test-darwin-arm64 (2026-07-15T22:29:10Z)\n\nMIT\n"
)

func TestNewVersionCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expected    string
		expectedErr string
	}{
		{"no arguments", nil, expectedLongVersion, ""},
		{"flag: --short", []string{"--short"}, expectedShortVersion, ""},
		{"shorthand flag: -s", []string{"-s"}, expectedShortVersion, ""},
		{"unknown flag", []string{"--unknown"}, "", "unknown flag"},
		{"unknown shorthand", []string{"-u"}, "", "unknown shorthand"},
		{"unknown command", []string{"unknown"}, "", "unknown command"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := new(bytes.Buffer)

			cmd := NewVersionCommand(testMetadata)
			cmd.SetOut(out)
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			if test.expectedErr == "" {
				assert.NilError(t, err)
				assert.Equal(t, out.String(), test.expected)
			} else {
				assert.ErrorContains(t, err, test.expectedErr)
			}
		})
	}
}
