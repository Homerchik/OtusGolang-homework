package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("Check read env vars from dir", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		_ = env
		expected := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		require.NoError(t, err, "Unexpected error during processing")
		require.Equal(t, expected, env)
	})

	t.Run("Check read env vars from bad dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/env/FOO")
		require.Error(t, err, "Unexpected error during processing")
	})

	t.Run("Check read env vars from unexisting dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/boo/boo")
		require.Error(t, err, "Unexpected error during processing")
	})

	t.Run("Check read env vars skip files with `=` in name", func(t *testing.T) {
		d, err := os.MkdirTemp("/tmp", "test")
		require.NoError(t, err)
		_, err = os.CreateTemp(d, "X=*")
		require.NoError(t, err)
		env, err := ReadDir(d)
		require.NoError(t, err, "Unexpected error during processing")
		require.Equal(t, Environment{}, env)
	})
}
