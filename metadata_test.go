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
