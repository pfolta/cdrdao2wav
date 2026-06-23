package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	ProjectName = "cdrdao2wav"
	Version     = "0.1.0"
	Copyright   = "Copyright (C) Peter Folta"
)

const (
	SectorSize      = 2352
	SampleRate      = 44100
	Channels        = 2
	BitsPerSample   = 16
	FramesPerSecond = 75
)

const (
	TocTrack       = "TRACK AUDIO"
	TocFile        = "FILE"
	TocCdTextTitle = "TITLE"
)

type Track struct {
	Number       int
	Title        string
	StartSector  int64
	LengthSector int64
	BinFile      string
}

var tocFileRegex = regexp.MustCompile(`^FILE\s+"(.+)"\s+(.+)\s+(.+)$`)
var msfRegex = regexp.MustCompile(`^(\d{2}):(0[0-9]|[1-5][0-9]):(0[0-9]|[1-6][0-9]|7[0-4])$`)

func writeWavHeader(w io.Writer, dataSize uint32) error {
	byteRate := uint32(SampleRate * Channels * BitsPerSample / 8)
	blockAlign := uint16(Channels * BitsPerSample / 8)

	header := struct {
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
	}{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + dataSize,
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   Channels,
		SampleRate:    SampleRate,
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: BitsPerSample,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: dataSize,
	}

	return binary.Write(w, binary.LittleEndian, &header)
}

func parseTocFile(tocFile string) ([]Track, error) {
	f, err := os.Open(tocFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var tracks []Track
	var current *Track
	trackNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, TocTrack) {
			trackNum++
			tracks = append(tracks, Track{Number: trackNum})
			current = &tracks[len(tracks)-1]
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, TocCdTextTitle) {
			current.Title = strings.Trim(line[6:], `"`)
			continue
		}

		if strings.HasPrefix(line, TocFile) {
			file := tocFileRegex.FindStringSubmatch(line)
			if file == nil {
				return nil, fmt.Errorf("Invalid FILE line: %s", line)
			}

			current.BinFile = file[1]

			start, err := msfToFrames(file[2])
			if err != nil {
				return nil, err
			}

			length, err := msfToFrames(file[3])
			if err != nil {
				return nil, err
			}

			current.StartSector = start
			current.LengthSector = length
		}
	}

	return tracks, scanner.Err()
}

func msfToFrames(offset string) (int64, error) {
	// Handle `00:00:00` formatted as `0`
	if offset == "0" {
		return 0, nil
	}

	msf := msfRegex.FindStringSubmatch(offset)
	if msf == nil {
		return 0, fmt.Errorf("Invalid MSF: %s", offset)
	}

	minutes, _ := strconv.Atoi(msf[1])
	seconds, _ := strconv.Atoi(msf[2])
	frames, _ := strconv.Atoi(msf[3])

	return int64((minutes*60+seconds)*FramesPerSecond + frames), nil
}

func validateTracks(tracks []Track, binSize int64) error {
	for _, track := range tracks {
		if track.LengthSector <= 0 {
			return fmt.Errorf("Track %d has invalid length", track.Number)
		}

		start := track.StartSector * SectorSize
		end := start + (track.LengthSector * SectorSize)

		if end > binSize {
			return fmt.Errorf("Track %d exceeds BIN size", track.Number)
		}
	}

	return nil
}

func extractTrack(binFileReader io.ReaderAt, track Track, wavFileWriter io.Writer) error {
	startByte := track.StartSector * SectorSize
	sizeBytes := track.LengthSector * SectorSize

	if sizeBytes > 0xFFFFFFFF {
		return fmt.Errorf("Track of size %d too large for WAV", sizeBytes)
	}

	if err := writeWavHeader(wavFileWriter, uint32(sizeBytes)); err != nil {
		return err
	}

	buf := make([]byte, 64*1024)
	remaining := sizeBytes

	for remaining > 0 {
		chunk := min(int64(len(buf)), remaining)

		n, err := binFileReader.ReadAt(buf[:chunk], startByte)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			swapEndian16(buf[:n])

			if _, err := wavFileWriter.Write(buf[:n]); err != nil {
				return err
			}

			startByte += int64(n)
			remaining -= int64(n)
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func swapEndian16(b []byte) {
	n := len(b) &^ 1
	for i := 0; i < n; i += 2 {
		b[i], b[i+1] = b[i+1], b[i]
	}
}

func sanitizeFilename(filename string) string {
	repl := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)

	return repl.Replace(strings.TrimSpace(filename))
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("%s version %s - %s\n", ProjectName, Version, Copyright)
			return
		}
	}

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <toc-file> <output-dir>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	tocFile := os.Args[1]
	outDir := os.Args[2]

	fmt.Printf("Parsing toc-file %q...\n", tocFile)
	tracks, err := parseTocFile(tocFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	if len(tracks) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: No tracks found in '%s'\n", tocFile)
		os.Exit(1)
	}

	// Resolve bin-file path relative to toc-file path
	tocDir := filepath.Dir(tocFile)

	for i := range tracks {
		if tracks[i].BinFile != "" {
			tracks[i].BinFile = filepath.Join(tocDir, tracks[i].BinFile)
		}
	}

	// Parse bin-file from first track
	binFile := tracks[0].BinFile
	if binFile == "" {
		fmt.Fprintf(os.Stderr, "ERROR: No bin-file found in '%s'\n", tocFile)
		os.Exit(1)
	}

	// Ensure toc-file only references one bin-file
	for _, track := range tracks {
		if track.BinFile != binFile {
			fmt.Fprintf(os.Stderr, "ERROR: Multiple bin-files are not supported\n")
			os.Exit(1)
		}
	}

	st, err := os.Stat(binFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	if err := validateTracks(tracks, st.Size()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	bin, err := os.Open(binFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	defer bin.Close()

	// Ensure outDir exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	for _, track := range tracks {
		wavFilename := fmt.Sprintf("%02d", track.Number)

		// Append track title from CD-TEXT if available
		if title := sanitizeFilename(track.Title); title != "" {
			wavFilename += " " + title
		}

		wavFilename += ".wav"

		wavFilePath := filepath.Join(outDir, wavFilename)

		wavFile, err := os.Create(wavFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Track %d: %v\n", track.Number, err)
			os.Exit(1)
		}
		defer wavFile.Close()

		fmt.Printf("Extracting track %d: %s...\n", track.Number, wavFilePath)

		if err := extractTrack(bin, track, wavFile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Track %d: %v\n", track.Number, err)
			os.Exit(1)
		}
	}

	fmt.Println("Done.")
}
