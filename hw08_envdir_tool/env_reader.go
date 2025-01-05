package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const delimiter byte = 0x0A

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readEnvValue(file *os.File) (*EnvValue, error) {
	data, err := bufio.NewReader(file).ReadBytes(delimiter)
	if err != nil && err != io.EOF {
		return nil, err
	}
	data = bytes.ReplaceAll(data, []byte{0x00}, []byte{delimiter})
	value := strings.TrimRight(string(data), " \t\n")
	return &EnvValue{Value: value, NeedRemove: false}, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	env := make(map[string]EnvValue)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}
		in, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		stat, err := in.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Size() == 0 {
			env[file.Name()] = EnvValue{Value: "", NeedRemove: true}
		} else {
			if envVar, err := readEnvValue(in); err == nil {
				env[file.Name()] = *envVar
			} else {
				return nil, err
			}
		}
	}
	return env, nil
}
