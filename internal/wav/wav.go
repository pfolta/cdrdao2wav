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
	ErrWriteIncomplete = errors.New("wav: write incomplete")
	ErrWriteOverflow   = errors.New("wav: write overflow")
	ErrCloseFile       = errors.New("wav: close file")
)

type Writer struct {
	file         *os.File
	dataSize     uint32
	bytesWritten int
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

	w := &Writer{
		file:     file,
		dataSize: dataSize,
	}

	if err := w.writeHeader(); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("%w: %w", ErrWriteHeader, err)
	}

	return w, nil
}

// Write writes PCM data to the WAV file.
// It returns the number of bytes written and an error, if any.
func (w *Writer) Write(p []byte) (int, error) {
	if w.closed {
		err := fmt.Errorf("%w: write after close", ErrWriteData)
		return 0, err
	}

	if w.bytesWritten+len(p) > int(w.dataSize) {
		err := fmt.Errorf(
			"%w: %d bytes written, attempting to write %d bytes, exceeds %d bytes",
			ErrWriteOverflow,
			w.bytesWritten,
			len(p),
			w.dataSize,
		)
		return 0, err
	}

	bytesWritten, err := w.file.Write(p)
	w.bytesWritten += bytesWritten

	if err != nil {
		return bytesWritten, fmt.Errorf("%w: %w", ErrWriteData, err)
	}

	return bytesWritten, nil
}

// Close closes the WAV file.
func (w *Writer) Close() error {
	if w.closed {
		return nil
	}

	err := w.file.Close()
	w.closed = true
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCloseFile, err)
	}

	if w.bytesWritten != int(w.dataSize) {
		return fmt.Errorf(
			"%w: wrote %d bytes, expected %d bytes",
			ErrWriteIncomplete,
			w.bytesWritten,
			w.dataSize,
		)
	}

	return nil
}

func (w *Writer) writeHeader() error {
	wavHeader := header{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + w.dataSize,
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
		Subchunk2Size: w.dataSize,
	}

	return binary.Write(w.file, binary.LittleEndian, &wavHeader)
}
