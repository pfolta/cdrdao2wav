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

package wav

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestWrite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		bytesWritten, err := writer.Write([]byte{1, 2, 3, 4})
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(bytesWritten, 4))

		bytesWritten, err = writer.Write([]byte{5, 6, 7, 8})
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(bytesWritten, 4))

		err = writer.Close()
		assert.NilError(t, err)

		fileInfo, err := os.Stat(path)
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(fileInfo.Size(), int64(binary.Size(header{})+int(dataSize))))

		data, err := os.ReadFile(path)
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(string(data[0:4]), "RIFF"))
		assert.Assert(t, is.Equal(string(data[8:12]), "WAVE"))
		assert.Assert(t, is.Equal(binary.LittleEndian.Uint32(data[24:28]), SampleRate))
		assert.Assert(t, is.Equal(binary.LittleEndian.Uint16(data[22:24]), Channels))
		assert.Assert(t, is.Equal(binary.LittleEndian.Uint16(data[34:36]), BitsPerSample))
		assert.Assert(t, is.Equal(binary.LittleEndian.Uint32(data[40:44]), dataSize))
		assert.Assert(t, is.DeepEqual(data[44:52], []byte{1, 2, 3, 4, 5, 6, 7, 8}))
	})

	t.Run("overflow", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		bytesWritten, err := writer.Write([]byte{1, 2, 3, 4})
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(bytesWritten, 4))

		bytesWritten, err = writer.Write([]byte{5, 6, 7, 8, 9})
		assert.Assert(t, is.ErrorIs(err, ErrWriteOverflow))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("%d bytes written", 4)))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("write %d bytes", 5)))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("exceeds %d bytes", dataSize)))
		assert.Assert(t, is.Equal(bytesWritten, 0))
	})

	t.Run("after close", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		_, err = writer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
		assert.NilError(t, err)

		err = writer.Close()
		assert.NilError(t, err)

		_, err = writer.Write([]byte{9})
		assert.Assert(t, is.ErrorIs(err, ErrWriteData))
		assert.Assert(t, is.ErrorContains(err, "write after close"))
	})
}

func TestClose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		_, err = writer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
		assert.NilError(t, err)

		assert.Assert(t, is.Equal(writer.closed, false))

		err = writer.Close()
		assert.NilError(t, err)

		assert.Assert(t, is.Equal(writer.closed, true))
	})

	t.Run("idempotent", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		_, err = writer.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
		assert.NilError(t, err)

		assert.Assert(t, is.Equal(writer.closed, false))

		for range 2 {
			err = writer.Close()
			assert.NilError(t, err)

			assert.Assert(t, is.Equal(writer.closed, true))
		}
	})

	t.Run("incomplete data", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.wav")
		dataSize := uint32(8)

		writer, err := NewWriter(path, dataSize)
		assert.NilError(t, err)

		bytesWritten, err := writer.Write([]byte{1, 2, 3, 4}) // incomplete
		assert.NilError(t, err)

		assert.Assert(t, is.Equal(writer.closed, false))

		err = writer.Close()
		assert.Assert(t, is.ErrorIs(err, ErrWriteIncomplete))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("wrote %d bytes", bytesWritten)))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("expected %d bytes", dataSize)))
	})
}
