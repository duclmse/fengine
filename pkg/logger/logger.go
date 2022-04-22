package logger

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"io"
	"time"
)

var _ Logger = (*logger)(nil)

// Logger specifies logging API.
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Elapse(format string, args ...interface{}) func(time.Time, *error)
	Struct(interface{})
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

type logger struct {
	kitLogger log.Logger
	level     Level
}

func (l logger) Debug(format string, args ...interface{}) {
	if Debug.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Debug.String(), "message", fmt.Sprintf(format, args...))
	}
}

func (l logger) Info(format string, args ...interface{}) {
	if Info.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Info.String(), "message", fmt.Sprintf(format, args...))
	}
}

func (l logger) Warn(format string, args ...interface{}) {
	if Warn.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Warn.String(), "message", fmt.Sprintf(format, args...))
	}
}

func (l logger) Error(format string, args ...interface{}) {
	if Error.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Error.String(), "message", fmt.Sprintf(format, args...))
	}
}

func (l logger) Fatalf(format string, args ...interface{}) {
	if Error.isAllowed(l.level) {
		msg := fmt.Sprintf(format, args...)
		_ = l.kitLogger.Log("level", Error.String(), "message", msg)
		panic(msg)
	}
}

func (l logger) Elapse(format string, args ...interface{}) func(time.Time, *error) {
	return func(begin time.Time, err *error) {
		if *err != nil {
			args = append(args, time.Since(begin), (*err).Error())
			_ = l.kitLogger.Log("level", Info.String(), "message", fmt.Sprintf(format+" took %s with error: %s.", args...))
			return
		}
		args = append(args, time.Since(begin))
		_ = l.kitLogger.Log("level", Info.String(), "message", fmt.Sprintf(format+" took %s without error.", args...))
	}
}

func (l logger) Struct(arg interface{}) {
	if Info.isAllowed(l.level) {
		_ = l.kitLogger.Log("level", Info.String(), "message", arg)
	}
}
