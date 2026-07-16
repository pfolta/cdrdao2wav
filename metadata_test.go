package cdrdao2audio

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetMetadata(t *testing.T) {
	originalVersion := version
	originalBuildDate := buildDate

	defer func() {
		version = originalVersion
		buildDate = originalBuildDate
	}()

	// Simulate build-time metadata injection.
	version = "v1.4.2-test"
	buildDate = "2026-07-15T22:29:10Z"

	metadata := GetMetadata()

	assert.Equal(t, metadata.Name, programName)
	assert.Equal(t, metadata.Version, "v1.4.2-test")
	assert.Equal(t, metadata.BuildDate, "2026-07-15T22:29:10Z")
	assert.Assert(t, metadata.License != "")
	assert.Assert(t, metadata.OS != "")
	assert.Assert(t, metadata.Arch != "")
}
