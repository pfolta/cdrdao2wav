package msf

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	MaxMinutes = 99
	MaxSeconds = 59
	MaxFrames  = 74

	FramesPerSecond = 75
	MaxLBA          = MaxMinutes*60*FramesPerSecond + MaxSeconds*FramesPerSecond + MaxFrames

	SectorBytes = 2352
)

var msfRegex = regexp.MustCompile(`^(\d{2}):(0[0-9]|[1-5][0-9]):(0[0-9]|[1-6][0-9]|7[0-4])$`)

// ErrorInvalidMSF indicates an invalid [MSF] value.
var ErrInvalidMSF = errors.New("invalid MSF")

// MSF represents an audio CD timestamp in minutes:seconds:frames (MM:SS:FF)
// stored as a Logical Block Address (LBA).
type MSF struct {
	lba uint32
}

// New returns a [MSF] for the specified Logical Block Address (LBA).
func New(lba uint32) (MSF, error) {
	if lba > MaxLBA {
		err := fmt.Errorf("%w: %d > %d", ErrInvalidMSF, lba, MaxLBA)
		return MSF{}, err
	}

	return MSF{lba}, nil
}

// Must is like [New] but panics if the specified LBA is invalid.
// It is intended for tests and trusted inputs.
func Must(lba uint32) MSF {
	msf, err := New(lba)
	if err != nil {
		panic(err)
	}

	return msf
}

// Parse parses a MM:SS:FF timestamp and returns the corresponding [MSF].
func Parse(str string) (MSF, error) {
	// Handle `00:00:00` formatted as `0`
	if str == "0" {
		return MSF{0}, nil
	}

	msfParts := msfRegex.FindStringSubmatch(str)
	if msfParts == nil {
		err := fmt.Errorf("%w: %s", ErrInvalidMSF, str)
		return MSF{}, err
	}

	m, _ := strconv.ParseUint(msfParts[1], 10, 32)
	s, _ := strconv.ParseUint(msfParts[2], 10, 32)
	f, _ := strconv.ParseUint(msfParts[3], 10, 32)

	return New(uint32(m*60*FramesPerSecond + s*FramesPerSecond + f))
}

// MustParse is like [Parse] but panics if the string cannot be parsed.
// It is intended for tests and trusted inputs.
func MustParse(str string) MSF {
	msf, err := Parse(str)
	if err != nil {
		panic(err)
	}

	return msf
}

// LBA returns the Logical Block Address of this [MSF] (total number of frames).
func (msf MSF) LBA() uint32 {
	return msf.lba
}

// SectorBytes returns the total size of this [MSF] in bytes.
func (msf MSF) SectorBytes() uint32 {
	return msf.lba * SectorBytes
}

// Add returns the sum of msf and other as a new [MSF] value.
// It returns an error if the resulting value would exceed the maximum MSF.
func (msf MSF) Add(other MSF) (MSF, error) {
	if msf.lba+other.lba > MaxLBA {
		err := fmt.Errorf("%w: %d + %d > %d", ErrInvalidMSF, msf.lba, other.lba, MaxLBA)
		return MSF{}, err
	}

	return MSF{msf.lba + other.lba}, nil
}

// Sub returns the difference between msf and other as a new [MSF] value.
// It returns an error if the resulting value would be negative.
func (msf MSF) Sub(other MSF) (MSF, error) {
	if other.lba > msf.lba {
		err := fmt.Errorf("%w: %d - %d < 0", ErrInvalidMSF, msf.lba, other.lba)
		return MSF{}, err
	}

	return MSF{msf.lba - other.lba}, nil
}

// String returns a MM:SS:FF timestamp representation of this [MSF].
func (msf MSF) String() string {
	m := msf.lba / FramesPerSecond / 60
	s := msf.lba / FramesPerSecond % 60
	f := msf.lba % FramesPerSecond
	return fmt.Sprintf("%02d:%02d:%02d", m, s, f)
}
