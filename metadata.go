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
	_ "embed"
	"runtime"
	"strings"
)

const programName = "cdrdao2audio"

// Metadata injected during build via -ldflags.
var (
	buildDate = "unknown"
	version   = "dev"
)

//go:embed LICENSE
var license string

// Metadata contains build and runtime information about the application.
type Metadata struct {
	// Name is the application name.
	Name string

	// Version is the application version injected at build time.
	Version string

	// BuildDate is the timestamp when the binary was built.
	BuildDate string

	// License contains the application license text.
	License string

	// OS is the target operating system.
	OS string

	// Arch is the target architecture.
	Arch string
}

// GetMetadata returns the application's build and runtime metadata.
func GetMetadata() Metadata {
	return Metadata{
		Name:      programName,
		Version:   version,
		BuildDate: buildDate,
		License:   strings.TrimSpace(license),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}
