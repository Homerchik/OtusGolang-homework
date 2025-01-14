package logger

import (
	"fmt"
	"strings"
	"time"
)

const (
	DEBUG = iota
	INFO
	ERROR
)

var LogLevel = map[string]int32{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"ERROR": ERROR,
}

type Logger struct {
	Level  int32
	Layout string
}

func New(level, timestampLayout string) *Logger {
	if timestampLayout == "" {
		timestampLayout = time.RFC3339
	}
	return &Logger{
		Level:  LogLevel[strings.ToUpper(level)],
		Layout: timestampLayout,
	}
}

func (l Logger) Debug(format string, a ...any) {
	if l.Level > DEBUG {
		return
	}
	l.PrintWithTime(format, "DEBUG", a...)
}

func (l Logger) Info(format string, a ...any) {
	if l.Level > INFO {
		return
	}
	l.PrintWithTime(format, "INFO", a...)
}

func (l Logger) Error(format string, a ...any) {
	l.PrintWithTime(format, "ERROR", a...)
}

func (l Logger) PrintWithTime(format, level string, a ...any) {
	format = strings.Join([]string{"%s", "%s", format, "\n"}, " ")
	timestamp := time.Now().UTC().Format(l.Layout)
	msg := fmt.Sprintf(format, a...)
	logLine := strings.Join([]string{timestamp, level, msg}, ": ")
	fmt.Println(logLine)
}
