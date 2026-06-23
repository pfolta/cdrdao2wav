package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
)

func TestWriteWavHeader(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := writeWavHeader(buffer, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := buffer.Bytes()

	if len(data) != 44 {
		t.Fatalf("header too short: got %d bytes", len(data))
	}

	// RIFF header
	if string(data[0:4]) != "RIFF" {
		t.Errorf("missing RIFF header")
	}

	// WAVE format
	if string(data[8:12]) != "WAVE" {
		t.Errorf("missing WAVE format")
	}

	if binary.LittleEndian.Uint32(data[24:28]) != SampleRate {
		t.Errorf("missing or incorrect sample rate")
	}

	if binary.LittleEndian.Uint16(data[22:24]) != Channels {
		t.Errorf("missing or incorrect number of channels")
	}

	if binary.LittleEndian.Uint16(data[34:36]) != BitsPerSample {
		t.Errorf("missing or incorrect bit per sample")
	}

	// data chunk
	if string(data[36:40]) != "data" {
		t.Errorf("missing data chunk")
	}
}

func TestParseTocFile(t *testing.T) {
	tests := []struct {
		in         string
		want       []Track
		wantErr    bool
		wantErrMsg string
	}{
		{
			"with-cd-text.toc",
			[]Track{
				{1, "The Fate of Ophelia", 0, 16955, "with-cd-text.bin"},
				{2, "Elizabeth Taylor", 16955, 15622, "with-cd-text.bin"},
				{3, "Opalite", 32577, 17652, "with-cd-text.bin"},
				{4, "Father Figure", 50229, 15958, "with-cd-text.bin"},
				{5, "Eldest Daughter", 66187, 18477, "with-cd-text.bin"},
				{6, "Ruin The Friendship", 84664, 16542, "with-cd-text.bin"},
				{7, "Actually Romantic", 101206, 12277, "with-cd-text.bin"},
				{8, "Wi$h Li$t", 113483, 15552, "with-cd-text.bin"},
				{9, "Wood", 129035, 11289, "with-cd-text.bin"},
				{10, "CANCELLED!", 140324, 15862, "with-cd-text.bin"},
				{11, "Honey", 156186, 13615, "with-cd-text.bin"},
				{12, "The Life of a Showgirl (Feat. Sabrina Carpenter)", 169801, 18133, "with-cd-text.bin"},
			},
			false,
			"",
		},
		{
			"without-cd-text.toc",
			[]Track{
				{1, "", 0, 17172, "without-cd-text.bin"},
				{2, "", 17172, 21979, "without-cd-text.bin"},
				{3, "", 39151, 15285, "without-cd-text.bin"},
				{4, "", 54436, 19592, "without-cd-text.bin"},
				{5, "", 74028, 19723, "without-cd-text.bin"},
				{6, "", 93751, 25532, "without-cd-text.bin"},
				{7, "", 119283, 15809, "without-cd-text.bin"},
				{8, "", 135092, 16160, "without-cd-text.bin"},
				{9, "", 151252, 19077, "without-cd-text.bin"},
				{10, "", 170329, 25057, "without-cd-text.bin"},
				{11, "", 195386, 11722, "without-cd-text.bin"},
				{12, "", 207108, 20787, "without-cd-text.bin"},
				{13, "", 227895, 16350, "without-cd-text.bin"},
				{14, "", 244245, 18416, "without-cd-text.bin"},
				{15, "", 262661, 14767, "without-cd-text.bin"},
				{16, "", 277428, 16250, "without-cd-text.bin"},
				{17, "", 293678, 16860, "without-cd-text.bin"},
			},
			false,
			"",
		},
		{"empty.toc", []Track{}, false, ""},

		{"invalid-msf.toc", []Track{}, true, "Invalid MSF"},
		{"invalid-file-line.toc", []Track{}, true, "Invalid FILE line"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseTocFile("testdata/" + tt.in)

			if tt.wantErr {
				if err == nil {
					t.Errorf("%q: got %v, expected error", tt.in, got)
					return
				}

				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("%q: got %q, want %q", tt.in, err.Error(), tt.wantErrMsg)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("%q: got %d tracks, want %d tracks", tt.in, len(got), len(tt.want))
			}

			for i := range got {
				if got[i].Number != tt.want[i].Number {
					t.Errorf("%v: got track number %d, want track number %d", tt.want[i], got[i].Number, tt.want[i].Number)
				}

				if got[i].Title != tt.want[i].Title {
					t.Errorf("%v: got track title %q, want track title %q", tt.want[i], got[i].Title, tt.want[i].Title)
				}

				if got[i].StartSector != tt.want[i].StartSector {
					t.Errorf("%v: got track start sector %d, want track start sector %d", tt.want[i], got[i].StartSector, tt.want[i].StartSector)
				}

				if got[i].LengthSector != tt.want[i].LengthSector {
					t.Errorf("%v: got track sector length %d, want track sector length %d", tt.want[i], got[i].LengthSector, tt.want[i].LengthSector)
				}

				if got[i].BinFile != tt.want[i].BinFile {
					t.Errorf("%v: got track bin-file %q, want track bin-file %q", tt.want[i], got[i].BinFile, tt.want[i].BinFile)
				}
			}
		})
	}
}

func TestMsfToFrames(t *testing.T) {
	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"0", 0, false},
		{"00:00:00", 0, false},
		{"00:00:01", 1, false},
		{"00:01:00", 75, false},
		{"01:00:00", 75 * 60, false},
		{"02:10:30", (2*60+10)*75 + 30, false},

		{"00:00", 0, true},
		{"00:00:00:00", 0, true},
		{"foobar", 0, true},
		{"aa:00:00", 0, true},
		{"00:bb:00", 0, true},
		{"00:00:cc", 0, true},
		{"01:02:75", 0, true},
		{"01:61:00", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := msfToFrames(tt.in)

			if tt.wantErr {
				if err == nil {
					t.Errorf("%q: got %d, expected error", tt.in, got)
					return
				}

				wantErr := "Invalid MSF: " + tt.in
				if err.Error() != wantErr {
					t.Errorf("%q: got %q, want %q", tt.in, err.Error(), wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("%q: unexpected error: %v", tt.in, err)
			}

			if got != tt.want {
				t.Errorf("%q: got %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateTracks(t *testing.T) {
	const binSize = int64(SectorSize * 100)

	tests := []struct {
		name    string
		tracks  []Track
		wantErr bool
	}{
		{
			name: "valid track",
			tracks: []Track{
				{Number: 1, StartSector: 0, LengthSector: 10},
			},
			wantErr: false,
		},
		{
			name: "multiple valid tracks",
			tracks: []Track{
				{Number: 1, StartSector: 0, LengthSector: 10},
				{Number: 2, StartSector: 10, LengthSector: 20},
			},
			wantErr: false,
		},
		{
			name: "zero length",
			tracks: []Track{
				{Number: 1, StartSector: 0, LengthSector: 0},
			},
			wantErr: true,
		},
		{
			name: "negative length",
			tracks: []Track{
				{Number: 1, StartSector: 0, LengthSector: -1},
			},
			wantErr: true,
		},
		{
			name: "exceeds bin-file size",
			tracks: []Track{
				{Number: 1, StartSector: 90, LengthSector: 20},
			},
			wantErr: true,
		},
		{
			name: "exact bin-file boundary is valid",
			tracks: []Track{
				{Number: 1, StartSector: 90, LengthSector: 10},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTracks(tt.tracks, binSize)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestExtractTrack(t *testing.T) {
	binData := make([]byte, SectorSize*2)
	for i := range binData {
		binData[i] = byte(i)
	}

	bin := bytes.NewReader(binData)

	var out bytes.Buffer

	track := Track{
		StartSector:  0,
		LengthSector: 2,
	}

	err := extractTrack(bin, track, &out)
	if err != nil {
		t.Fatalf("extractTrack failed: %v", err)
	}

	data := out.Bytes()

	if len(data) < 44 {
		t.Fatalf("output too short")
	}

	if string(data[0:4]) != "RIFF" {
		t.Errorf("missing RIFF")
	}

	if string(data[8:12]) != "WAVE" {
		t.Errorf("missing WAVE")
	}

	if string(data[36:40]) != "data" {
		t.Errorf("missing data chunk")
	}

	expected := SectorSize * 2
	if len(data)-44 != expected {
		t.Errorf("data size mismatch: got %d want %d", len(data)-44, expected)
	}
}

func TestSwapEndian16(t *testing.T) {
	tests := []struct {
		in   []byte
		want []byte
	}{
		{
			in:   []byte{0x01, 0x02},
			want: []byte{0x02, 0x01},
		},
		{
			in:   []byte{0x01, 0x02, 0x03, 0x04},
			want: []byte{0x02, 0x01, 0x04, 0x03},
		},
		{
			// odd length: last byte should remain untouched
			in:   []byte{0x01, 0x02, 0x03},
			want: []byte{0x02, 0x01, 0x03},
		},
		{
			// empty slice
			in:   []byte{},
			want: []byte{},
		},
		{
			// single byte
			in:   []byte{0xFF},
			want: []byte{0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.in), func(t *testing.T) {
			// copy input so we don't mutate test case data
			got := append([]byte(nil), tt.in...)

			swapEndian16(got)

			if !bytes.Equal(got, tt.want) {
				t.Errorf("%#v: got %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Track 01", "Track 01"},
		{"Track:One", "Track-One"},
		{"Bad/Name", "Bad-Name"},
		{"Weird|Name?", "WeirdName"},
		{"  spaced  ", "spaced"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sanitizeFilename(tt.in)

			if got != tt.want {
				t.Errorf("%q: got %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
