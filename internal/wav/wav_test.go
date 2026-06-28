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
		assert.Assert(t, is.Equal(bytesWritten, uint32(4)))

		bytesWritten, err = writer.Write([]byte{5, 6, 7, 8})
		assert.NilError(t, err)
		assert.Assert(t, is.Equal(bytesWritten, uint32(4)))

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
		assert.Assert(t, is.Equal(bytesWritten, uint32(4)))

		bytesWritten, err = writer.Write([]byte{5, 6, 7, 8, 9})
		assert.Assert(t, is.ErrorIs(err, ErrDataOverflow))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("%d bytes written", 4)))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("write %d bytes", 5)))
		assert.Assert(t, is.ErrorContains(err, fmt.Sprintf("exceeds %d bytes", dataSize)))
		assert.Assert(t, is.Equal(bytesWritten, uint32(0)))
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
