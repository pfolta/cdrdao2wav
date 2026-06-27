package msf

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestNew(t *testing.T) {
	tests := []uint32{
		0,
		1000,
		449_999,
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("New(%d)", test), func(t *testing.T) {
			actual, err := New(test)
			assert.Assert(t, is.Nil(err))
			assert.Equal(t, actual.TotalFrames(), test)
		})
	}
}

func TestNewErrors(t *testing.T) {
	tests := []uint32{
		450_000,
		1_000_000,
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("New(%d) errors", test), func(t *testing.T) {
			actual, err := New(test)
			assert.Assert(t, is.ErrorIs(err, ErrInvalidMSF))
			assert.Assert(t, is.Equal(actual, MSF{}))
		})
	}
}

func TestMust(t *testing.T) {
	t.Run("Must(1000) does not panic", func(t *testing.T) {
		assert.Assert(t, is.Equal(Must(1_000), MSF{1_000}))
	})

	t.Run("Must(1000000) panics", func(t *testing.T) {
		assert.Assert(t, is.Panics(func() { Must(1_000_000) }))
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		in       string
		expected MSF
	}{
		{"0", MSF{0}},
		{"00:00:00", MSF{0}},
		{"00:00:01", MSF{1}},
		{"00:01:00", MSF{75}},
		{"01:00:00", MSF{4_500}},
		{"02:10:30", MSF{9_780}},
		{"99:59:74", MSF{449_999}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Parse(%s)", test.in), func(t *testing.T) {
			actual, err := Parse(test.in)
			assert.Assert(t, is.Nil(err))
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"",
		"00:00",
		"00:00:00:00",
		"foobar",
		"aa:00:00",
		"00:bb:00",
		"00:00:cc",
		"01:02:75",
		"01:61:00",
		"-01:02:03",
		"01:-02:03",
		"01:02:-03",
		"-01:-02:-03",
		"123:123:123",
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Parse(%q) errors", test), func(t *testing.T) {
			actual, err := Parse(test)
			assert.Assert(t, is.ErrorIs(err, ErrInvalidMSF))
			assert.Assert(t, is.Equal(actual, MSF{}))
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("MustParse(\"01:00:00\") does not panic", func(t *testing.T) {
		assert.Assert(t, is.Equal(MustParse("01:00:00"), MSF{4_500}))
	})

	t.Run("MustParse(\"invalid\") panics", func(t *testing.T) {
		assert.Assert(t, is.Panics(func() { MustParse("invalid") }))
	})
}

func TestTotalFrames(t *testing.T) {
	tests := []struct {
		in       MSF
		expected uint32
	}{
		{MSF{}, 0},
		{MSF{0}, 0},
		{MSF{1_000}, 1_000},
		{MSF{449_999}, 449_999},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.TotalFrames()", test.in.TotalFrames()), func(t *testing.T) {
			actual := test.in.TotalFrames()
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestSectorBytes(t *testing.T) {
	tests := []struct {
		in       MSF
		expected uint32
	}{
		{MSF{}, 0},
		{MSF{0}, 0},
		{MSF{1_000}, 2_352_000},
		{MSF{449_999}, 1_058_397_648},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.SectorBytes()", test.in.TotalFrames()), func(t *testing.T) {
			actual := test.in.SectorBytes()
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		in       MSF
		expected string
	}{
		{MSF{}, "00:00:00"},
		{MSF{0}, "00:00:00"},
		{MSF{1000}, "00:13:25"},
		{MSF{449_999}, "99:59:74"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.String()", test.in.TotalFrames()), func(t *testing.T) {
			actual := test.in.String()
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}
