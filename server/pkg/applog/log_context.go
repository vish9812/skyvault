package applog

import (
	"time"

	"github.com/rs/zerolog"
)

var _ LogContext = (*zeroLogContext)(nil)

// LogContext is for creating a NEW logger with preset fields.
// Use the Logger.With() method to create a new LogContext.
// Populate the LogContext with the available methods.
// Use the LogContext.Logger() method to get a new Logger with the populated fields.
// Use the LogEvent.Msg() method to actually write the log.
type LogContext interface {
	Str(key, val string) LogContext
	Int(key string, val int) LogContext
	Int64(key string, val int64) LogContext
	Float64(key string, val float64) LogContext
	Bool(key string, val bool) LogContext
	Time(key string, val time.Time) LogContext
	Any(key string, val any) LogContext
	Logger() Logger
}

type zeroLogContext struct {
	ctx zerolog.Context
}

// Implementation of LogContext interface
func (c *zeroLogContext) Str(key, val string) LogContext {
	return &zeroLogContext{ctx: c.ctx.Str(key, val)}
}

func (c *zeroLogContext) Int(key string, val int) LogContext {
	return &zeroLogContext{ctx: c.ctx.Int(key, val)}
}

func (c *zeroLogContext) Int64(key string, val int64) LogContext {
	return &zeroLogContext{ctx: c.ctx.Int64(key, val)}
}

func (c *zeroLogContext) Float64(key string, val float64) LogContext {
	return &zeroLogContext{ctx: c.ctx.Float64(key, val)}
}

func (c *zeroLogContext) Bool(key string, val bool) LogContext {
	return &zeroLogContext{ctx: c.ctx.Bool(key, val)}
}

func (c *zeroLogContext) Time(key string, val time.Time) LogContext {
	return &zeroLogContext{ctx: c.ctx.Time(key, val)}
}

func (c *zeroLogContext) Any(key string, val any) LogContext {
	return &zeroLogContext{ctx: c.ctx.Interface(key, val)}
}

func (c *zeroLogContext) Logger() Logger {
	return &zeroLogger{log: c.ctx.Logger()}
}
