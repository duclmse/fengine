package logger

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"io"
	"time"
)

var _ Logger = (*logger)(nil)

type logger struct {
	kitLogger log.Logger
	level     Level
}

// Logger specifies logging API.
type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Debug_(format string, args ...interface{})
	Info_(format string, args ...interface{})
	Warn_(format string, args ...interface{})
	Error_(format string, args ...interface{})
}

// New returns wrapped go kit logger.
func New(out io.Writer, levelText string) (Logger, error) {
	var level Level
	err := level.UnmarshalText(levelText)
	if err != nil {
		return nil, fmt.Errorf(`{"level":"error","message":"%s: %s","ts":"%s"}`, err, levelText, time.RFC3339Nano)
	}
	l := log.NewJSONLogger(log.NewSyncWriter(out))
	l = log.With(l, "ts", log.DefaultTimestampUTC)
	return &logger{l, level}, err
}

func (l logger) Info(msg string) {
	if Info.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Info.String(), "message", msg)
	}
}

func (l logger) Debug(msg string) {
	if Debug.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Debug.String(), "message", msg)
	}
}

func (l logger) Warn(msg string) {
	if Warn.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Warn.String(), "message", msg)
	}
}

func (l logger) Error(msg string) {
	if Error.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Error.String(), "message", msg)
	}
}

func (l logger) Debug_(msg string, opts ...interface{}) {
	if Debug.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Debug.String(), "message", fmt.Sprintf(msg, opts))
	}
}

func (l logger) Warn_(msg string, opts ...interface{}) {
	if Warn.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Warn.String(), "message", fmt.Sprintf(msg, opts))
	}
}

func (l logger) Error_(msg string, opts ...interface{}) {
	if Error.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Error.String(), "message", fmt.Sprintf(msg, opts))
	}
}

func (l logger) Info_(msg string, opts ...interface{}) {
	if Info.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Info.String(), "message", fmt.Sprintf(msg, opts))
	}
}
