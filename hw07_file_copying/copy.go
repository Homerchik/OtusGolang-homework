package main

import (
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	bulkSize int64 = 256
)

func checkInFile(file *os.File, offset int64) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	size := stat.Size()
	if offset > size {
		return 0, ErrOffsetExceedsFileSize
	}
	if size == 0 {
		return 0, ErrUnsupportedFile
	}
	return size, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	size, err := checkInFile(fromFile, offset)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	bar := progressbar.NewOptions(100)

	bytesRead := int64(0)
	buf := make([]byte, bulkSize)
	if limit == 0 {
		limit = size
	}
	total := min(limit, max(0, size-offset))
	for bytesRead < limit {
		n, err := fromFile.ReadAt(buf, offset+bytesRead)
		toWrite := min(limit-bytesRead, bulkSize, int64(n))
		bytesRead += toWrite
		bar.Set(int(100 * bytesRead / total))
		toFile.Write(buf[:toWrite])
		if err == io.EOF {
			return nil
		}
	}
	return nil
}
