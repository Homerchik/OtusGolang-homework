package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

const chunk int64 = 256

// var replacer = strings.NewReplacer("\x00", "\n")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readEnvValue(file *os.File) (*EnvValue, error) {
	buf := make([]byte, chunk)
	var line []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n > 0 {
			line = append(line, buf...)
		}
		if strings.Contains(string(line), "\n") || err == io.EOF {
			line := bytes.ReplaceAll(
				bytes.Split(line, []byte{0x0A})[0],
				[]byte{0x00},
				[]byte{0x0A},
			)
			value := strings.TrimRight(string(line), " \t\n")
			return &EnvValue{Value: value, NeedRemove: false}, nil
		}
	}
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
		in, err := os.Open(dir + "/" + file.Name())
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
