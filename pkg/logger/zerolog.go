package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

type zerologLogger struct {
	wrappedLogger zerolog.Logger
}

func (c *zerologContext) formatMessage(message string) string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d %s", file[strings.LastIndex(file, "/")+1:], line, message)
	}
	return message
}

func newConsoleZerolog(level int) zerologLogger {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
	}

	writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	l := zerolog.New(writer).With().Timestamp().Stack().Logger()
	l.Level(levelSwitch(level))

	return zerologLogger{l}
}

func newTestZerolog(t *testing.T) zerologLogger {
	l := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Caller().Logger()

	return zerologLogger{l}
}

func (l zerologLogger) Debug() LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.DebugLevel,
	}
}
func (l zerologLogger) Info() LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.InfoLevel,
	}
}

func (l zerologLogger) Warn() LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.WarnLevel,
	}
}
func (l zerologLogger) Error() LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.ErrorLevel,
	}
}
func (l zerologLogger) Fatal() LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.FatalLevel,
	}
}

func (l zerologLogger) Err(err error) LogContext {
	return &zerologContext{
		logger: l.wrappedLogger,
		level:  zerolog.ErrorLevel,
		err:    err,
	}
}

type zerologContext struct {
	logger zerolog.Logger
	level  zerolog.Level
	err    error
}

func (c *zerologContext) Msg(message string) {
	if c.err != nil {
		c.logger.WithLevel(c.level).Err(c.err).Msg(c.formatMessage(message))
		return
	}
	c.logger.WithLevel(c.level).Msg(message)
}

func (c *zerologContext) Msgf(format string, args ...interface{}) {
	if c.err != nil {
		c.logger.WithLevel(c.level).Err(c.err).Msg(c.formatMessage(fmt.Sprintf(format, args...)))
		return
	}
	c.logger.WithLevel(c.level).Msgf(format, args...)
}

func (c *zerologContext) Debug() LogContext {
	c.level = zerolog.DebugLevel
	return c
}
func (c *zerologContext) Info() LogContext {
	c.level = zerolog.InfoLevel
	return c
}
func (c *zerologContext) Warn() LogContext {
	c.level = zerolog.WarnLevel
	return c
}
func (c *zerologContext) Error() LogContext {
	c.level = zerolog.ErrorLevel
	return c
}
func (c *zerologContext) Fatal() LogContext {
	c.level = zerolog.FatalLevel
	return c
}
func (c *zerologContext) Err(err error) LogContext {
	c.err = err
	return c
}
