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
