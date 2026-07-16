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

func TestRunVersion(t *testing.T) {
	tests := []struct {
		name     string
		opts     versionOptions
		expected string
	}{
		{"default", versionOptions{}, expectedLongVersion},
		{"short version", versionOptions{short: true}, expectedShortVersion},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := new(bytes.Buffer)

			err := runVersion(out, testMetadata, test.opts)
			assert.NilError(t, err)
			assert.Equal(t, out.String(), test.expected)
		})
	}
}
