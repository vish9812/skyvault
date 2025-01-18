package applog

import (
	"context"
	"os"
	"skyvault/pkg/common"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var _ Logger = (*zeroLogger)(nil)

// Logger is the main logging interface
// Use the NewLogger() function to create a new Logger.
// Use the Logger.<Level> method to create a new LogEvent when you want to write a log.
// Use the Logger.With() method when you only want to populate the data, but not write the log immediately.
// Use the LogContext.Logger() method to get a new Logger with the populated fields.
type Logger interface {
	Debug() LogEvent
	Info() LogEvent
	Warn() LogEvent
	Error() LogEvent
	Fatal() LogEvent
	With() LogContext
}

type zeroLogger struct {
	log zerolog.Logger
}

// Config for creating a new logger
type Config struct {
	Level      string
	TimeFormat string
	Console    bool
}

// NewLogger creates a new Logger instance
func NewLogger(config *Config) Logger {
	if config == nil {
		config = &Config{
			Level:      "info",
			TimeFormat: time.RFC3339,
			Console:    true,
		}
	}

	// Configure zerolog
	var output zerolog.ConsoleWriter
	if config.Console {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: config.TimeFormat,
			NoColor:    false,
		}
	}

	zl := zerolog.New(output).
		Level(parseLevel(config.Level)).
		With().
		Timestamp().
		Logger()

	return &zeroLogger{log: zl}
}

// parseLevel converts string level to zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func GetLoggerFromContext(ctx context.Context) Logger {
	return ctx.Value(common.CtxKeyLogger).(Logger)
}

// Implementation of Logger interface
func (l *zeroLogger) Debug() LogEvent {
	return &zeroLogEvent{evt: l.log.Debug()}
}

func (l *zeroLogger) Info() LogEvent {
	return &zeroLogEvent{evt: l.log.Info()}
}

func (l *zeroLogger) Warn() LogEvent {
	return &zeroLogEvent{evt: l.log.Warn()}
}

func (l *zeroLogger) Error() LogEvent {
	return &zeroLogEvent{evt: l.log.Error()}
}

func (l *zeroLogger) Fatal() LogEvent {
	return &zeroLogEvent{evt: l.log.Fatal()}
}

func (l *zeroLogger) With() LogContext {
	return &zeroLogContext{ctx: l.log.With()}
}
