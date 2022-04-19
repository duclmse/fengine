package logger

import (
	"errors"
	"strings"
)

// Level represents severity level while logging.
type Level uint

const (
	Error Level = iota + 1
	Warn
	Info
	Debug
	Trace
)

var ErrInvalidLogLevel = errors.New("unrecognized log level")

var levels = map[Level]string{
	Error: "error",
	Warn:  "warn",
	Info:  "info",
	Debug: "debug",
	Trace: "trace",
}

func (lvl Level) String() string {
	return levels[lvl]
}

func (lvl Level) isAllowed(logLevel Level) bool {
	return lvl <= logLevel
}

func (lvl *Level) UnmarshalText(text string) error {
	switch strings.ToLower(text) {
	case "trace":
		*lvl = Trace
	case "debug":
		*lvl = Debug
	case "info":
		*lvl = Info
	case "warn":
		*lvl = Warn
	case "error":
		*lvl = Error
	default:
		return ErrInvalidLogLevel
	}
	return nil
}
