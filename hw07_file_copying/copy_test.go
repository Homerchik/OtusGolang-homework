package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareFiles(t *testing.T, fromPath, toPath string, offset int64) {
	t.Helper()
	fromFile, _ := os.Open(fromPath)
	defer fromFile.Close()
	toFile, _ := os.Open(toPath)
	defer toFile.Close()
	stat, _ := toFile.Stat()
	bytesToCompare := stat.Size()
	fromFile.Seek(offset, 0)
	fromBytes := make([]byte, bytesToCompare)
	toBytes := make([]byte, bytesToCompare)
	require.Equal(t, fromBytes, toBytes, "Check files for equality")
}

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("Error opening file without size", func(t *testing.T) {
		err := Copy("/dev/urandom", "/any/path", 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("Error having offset greater that filesize", func(t *testing.T) {
		err := Copy("./testdata/out_offset6000_limit1000.txt", "/any/path", 1000, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("Check full-file copy, limit 0", func(t *testing.T) {
		fromFileName := "./testdata/input.txt"
		file, err := os.CreateTemp("/tmp", "output_1.txt")
		require.NoError(t, err)
		defer file.Close()
		require.NoError(t, err)
		err = Copy(fromFileName, file.Name(), 0, 0)
		require.NoError(t, err)
		compareFiles(t, fromFileName, file.Name(), 0)
	})

	t.Run("Check partial-file copy, limit 100, offset 100", func(t *testing.T) {
		fromFileName := "./testdata/input.txt"
		file, err := os.CreateTemp("/tmp", "output_1.txt")
		require.NoError(t, err)
		defer file.Close()
		err = Copy(fromFileName, file.Name(), 100, 100)
		require.NoError(t, err)
		compareFiles(t, fromFileName, file.Name(), 100)
	})
}
