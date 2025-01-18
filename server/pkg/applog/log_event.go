package applog

import (
	"time"

	"github.com/rs/zerolog"
)

var _ LogEvent = (*zeroLogEvent)(nil)

// LogEvent is for building AND writing the log.
// Use the Logger.<Level> method to create a new LogEvent.
// Populate the LogEvent with the available methods.
// Use the Write() method to actually write the log.
type LogEvent interface {
	Str(key, val string) LogEvent
	Int(key string, val int) LogEvent
	Int64(key string, val int64) LogEvent
	Float64(key string, val float64) LogEvent
	Bool(key string, val bool) LogEvent
	Time(key string, val time.Time) LogEvent
	Any(key string, val any) LogEvent
	Err(err error) LogEvent
	Write(msg string)
}

type zeroLogEvent struct {
	evt *zerolog.Event
}


// Implementation of LogEvent interface
func (e *zeroLogEvent) Str(key, val string) LogEvent {
	return &zeroLogEvent{evt: e.evt.Str(key, val)}
}

func (e *zeroLogEvent) Int(key string, val int) LogEvent {
	return &zeroLogEvent{evt: e.evt.Int(key, val)}
}

func (e *zeroLogEvent) Int64(key string, val int64) LogEvent {
	return &zeroLogEvent{evt: e.evt.Int64(key, val)}
}

func (e *zeroLogEvent) Float64(key string, val float64) LogEvent {
	return &zeroLogEvent{evt: e.evt.Float64(key, val)}
}

func (e *zeroLogEvent) Bool(key string, val bool) LogEvent {
	return &zeroLogEvent{evt: e.evt.Bool(key, val)}
}

func (e *zeroLogEvent) Time(key string, val time.Time) LogEvent {
	return &zeroLogEvent{evt: e.evt.Time(key, val)}
}

func (e *zeroLogEvent) Any(key string, val any) LogEvent {
	return &zeroLogEvent{evt: e.evt.Interface(key, val)}
}

func (e *zeroLogEvent) Err(err error) LogEvent {
	return &zeroLogEvent{evt: e.evt.Err(err)}
}

func (e *zeroLogEvent) Write(msg string) {
	e.evt.Msg(msg)
}