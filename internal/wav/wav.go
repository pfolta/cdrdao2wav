package wav

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const (
	SampleRate    uint32 = 44100
	Channels      uint16 = 2
	BitsPerSample uint16 = 16
)

var (
	ErrCreateFile      = errors.New("wav: create file")
	ErrWriteHeader     = errors.New("wav: write header")
	ErrWriteData       = errors.New("wav: write data")
	ErrDataOverflow    = errors.New("wav: write overflow")
	ErrWriteIncomplete = errors.New("wav: incomplete data")
	ErrCloseFile       = errors.New("wav: close file")
)

type Writer struct {
	file         *os.File
	dataSize     uint32
	bytesWritten uint32
	closed       bool
}

type header struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

// NewWriter creates a WAV file and returns a streaming writer for PCM data.
func NewWriter(path string, dataSize uint32) (*Writer, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreateFile, err)
	}

	writer := &Writer{
		file:     file,
		dataSize: dataSize,
	}

	if err := writer.writeHeader(); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("%w: %w", ErrWriteHeader, err)
	}

	return writer, nil
}

// Write writes a slice of bytes to the WAV file.
// It returns the number of bytes written and an error, if any.
func (writer *Writer) Write(chunk []byte) (uint32, error) {
	if writer.closed {
		err := fmt.Errorf("%w: write after close", ErrWriteData)
		return 0, err
	}

	if writer.bytesWritten+uint32(len(chunk)) > writer.dataSize {
		err := fmt.Errorf(
			"%w: %d bytes written, attempting to write %d bytes, exceeds %d bytes",
			ErrDataOverflow,
			writer.bytesWritten,
			len(chunk),
			writer.dataSize,
		)
		return 0, err
	}

	bytesWritten, err := writer.file.Write(chunk)
	writer.bytesWritten += uint32(bytesWritten)

	if err != nil {
		return uint32(bytesWritten), fmt.Errorf("%w: %w", ErrWriteData, err)
	}

	return uint32(bytesWritten), nil
}

// Close closes the WAV file.
func (writer *Writer) Close() error {
	if writer.closed {
		return nil
	}

	err := writer.file.Close()
	writer.closed = true
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCloseFile, err)
	}

	if writer.bytesWritten != writer.dataSize {
		return fmt.Errorf(
			"%w: wrote %d bytes, expected %d bytes",
			ErrWriteIncomplete,
			writer.bytesWritten,
			writer.dataSize,
		)
	}

	return nil
}

func (writer *Writer) writeHeader() error {
	header := header{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + writer.dataSize,
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   Channels,
		SampleRate:    SampleRate,
		ByteRate:      SampleRate * uint32(Channels) * uint32(BitsPerSample) / 8,
		BlockAlign:    Channels * BitsPerSample / 8,
		BitsPerSample: BitsPerSample,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: writer.dataSize,
	}

	return binary.Write(writer.file, binary.LittleEndian, &header)
}
