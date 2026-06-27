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

	MaxTotalFrames = MaxMinutes*60*75 + MaxSeconds*75 + MaxFrames

	SectorSize = 2352
)

var msfRegex = regexp.MustCompile(`^(\d{2}):(\d{2}):(\d{2})$`)

// ErrorInvalidMSF indicates an invalid [MSF] value.
var ErrInvalidMSF = errors.New("invalid MSF")

// MSF represents an audio CD timestamp in minutes:seconds:frames (MM:SS:FF).
type MSF struct {
	totalFrames uint32
}

// New returns a [MSF] for the specified total number of frames.
func New(frames uint32) (MSF, error) {
	if frames > MaxTotalFrames {
		err := fmt.Errorf("%w: %d > %d", ErrInvalidMSF, frames, MaxTotalFrames)
		return MSF{}, err
	}

	return MSF{frames}, nil
}

// Must is like [New] but panics if frames is invalid.
// It is intended for tests and trusted inputs.
func Must(frames uint32) MSF {
	msf, err := New(frames)
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

	if m > MaxMinutes || s > MaxSeconds || f > MaxFrames {
		err := fmt.Errorf("%w: %s", ErrInvalidMSF, str)
		return MSF{}, err
	}

	return New(uint32(m*60*75 + s*75 + f))
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

// TotalFrames returns the total number of frames of this [MSF].
func (msf MSF) TotalFrames() uint32 {
	return msf.totalFrames
}

// SectorBytes returns the total size of this [MSF] in bytes.
func (msf MSF) SectorBytes() uint32 {
	return msf.totalFrames * SectorSize
}

// String returns a MM:SS:FF timestamp representation of this [MSF].
func (msf MSF) String() string {
	m := msf.totalFrames / 75 / 60
	s := msf.totalFrames / 75 % 60
	f := msf.totalFrames % 75
	return fmt.Sprintf("%02d:%02d:%02d", m, s, f)
}
