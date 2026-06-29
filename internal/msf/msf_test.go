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
			assert.NilError(t, err)
			assert.Equal(t, actual.LBA(), test)
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
		{"00:00:74", MSF{74}},
		{"00:01:00", MSF{75}},
		{"01:00:00", MSF{4_500}},
		{"02:10:30", MSF{9_780}},
		{"99:59:74", MSF{449_999}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Parse(%s)", test.in), func(t *testing.T) {
			actual, err := Parse(test.in)
			assert.NilError(t, err)
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
	t.Run("MustParse(\"01:02:03\") does not panic", func(t *testing.T) {
		assert.Assert(t, is.Equal(MustParse("01:02:03"), MSF{4_653}))
	})

	t.Run("MustParse(\"invalid\") panics", func(t *testing.T) {
		assert.Assert(t, is.Panics(func() { MustParse("invalid") }))
	})
}

func TestLBA(t *testing.T) {
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
		t.Run(fmt.Sprintf("MSF{%d}.LBA()", test.in.LBA()), func(t *testing.T) {
			actual := test.in.LBA()
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
		t.Run(fmt.Sprintf("MSF{%d}.SectorBytes()", test.in.LBA()), func(t *testing.T) {
			actual := test.in.SectorBytes()
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a        MSF
		b        MSF
		expected MSF
	}{
		{MSF{}, MSF{}, MSF{}},
		{MSF{0}, MSF{0}, MSF{0}},
		{MSF{1}, MSF{2}, MSF{3}},
		{MSF{1_000}, MSF{2_000}, MSF{3_000}},
		{MSF{449_998}, MSF{1}, MSF{449_999}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.Add(MSF{%d})", test.a.LBA(), test.b.LBA()), func(t *testing.T) {
			actual, err := test.a.Add(test.b)
			assert.NilError(t, err)
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestAddErrors(t *testing.T) {
	tests := []struct {
		a MSF
		b MSF
	}{
		{MSF{449_999}, MSF{1}},
		{MSF{300_000}, MSF{300_000}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.Add(MSF{%d}) errors", test.a.LBA(), test.b.LBA()), func(t *testing.T) {
			actual, err := test.a.Add(test.b)
			assert.Assert(t, is.ErrorIs(err, ErrInvalidMSF))
			assert.Assert(t, is.Equal(actual, MSF{}))
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		a        MSF
		b        MSF
		expected MSF
	}{
		{MSF{}, MSF{}, MSF{}},
		{MSF{0}, MSF{0}, MSF{0}},
		{MSF{3}, MSF{2}, MSF{1}},
		{MSF{2_000}, MSF{1_000}, MSF{1_000}},
		{MSF{449_999}, MSF{1}, MSF{449_998}},
		{MSF{449_999}, MSF{449_999}, MSF{0}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.Sub(MSF{%d})", test.a.LBA(), test.b.LBA()), func(t *testing.T) {
			actual, err := test.a.Sub(test.b)
			assert.NilError(t, err)
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}

func TestSubErrors(t *testing.T) {
	tests := []struct {
		a MSF
		b MSF
	}{
		{MSF{0}, MSF{1}},
		{MSF{100_000}, MSF{200_000}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("MSF{%d}.Sub(MSF{%d}) errors", test.a.LBA(), test.b.LBA()), func(t *testing.T) {
			actual, err := test.a.Sub(test.b)
			assert.Assert(t, is.ErrorIs(err, ErrInvalidMSF))
			assert.Assert(t, is.Equal(actual, MSF{}))
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
		t.Run(fmt.Sprintf("MSF{%d}.String()", test.in.LBA()), func(t *testing.T) {
			actual := test.in.String()
			assert.Assert(t, is.Equal(actual, test.expected))
		})
	}
}
