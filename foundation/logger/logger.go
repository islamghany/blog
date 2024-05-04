package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

// different types of level for logging.
type Level slog.Level

// all possible logging level that can happen in the app
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// ===========================================

// EventFunc the the func type that will be invoked when the current event type happens.
type EventFunc func(ctx context.Context, r Record)

type Events struct {
	Debug EventFunc
	Info  EventFunc
	Warn  EventFunc
	Error EventFunc
}

// ===========================================

// Record represent the structs that will be logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

// generate a new Reocrd the slog.Record struct.
func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())
	f := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value.Any()
		return true
	}
	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: atts,
	}
}

// TraceIDFunc represents a function that can return the trace id from the specified context.
type TracerIDFunc func(ctx context.Context) string

// Logger struct
type Logger struct {
	handler      slog.Handler
	tracerIDFunc TracerIDFunc
}

// New constructs a new log for application use.
func New(w io.Writer, minLevel Level, serviceName string, traceIDFunc TracerIDFunc) *Logger {
	return new(w, minLevel, serviceName, traceIDFunc, Events{})
}

// NewWithEvents constructs a new log for application use with events.
func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFunc TracerIDFunc, events Events) *Logger {
	return new(w, minLevel, serviceName, traceIDFunc, events)
}

func NewStandardLogger(log *Logger, minLevel Level) *log.Logger {
	return slog.NewLogLogger(log.handler, slog.Level(minLevel))
}

// Debug logs at LevelDebug with the given context.
func (log *Logger) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelDebug, 3, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (log *Logger) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack position.
func (log *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelInfo, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (log *Logger) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelError, 3, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (log *Logger) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelWarn, 3, msg, args...)
}

func (log *Logger) write(ctx context.Context, logLevel Level, caller int, msg string, args ...any) {
	slogLevel := slog.Level(logLevel)
	if !log.handler.Enabled(ctx, slogLevel) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), slogLevel, msg, pcs[0])

	if log.tracerIDFunc != nil {
		args = append(args, "trace_id", log.tracerIDFunc(ctx))
	}
	r.Add(args...)
	log.handler.Handle(ctx, r)
}

func new(w io.Writer, minLevel Level, serviceName string, tracerIDFunc TracerIDFunc, events Events) *Logger {
	// Convert the file name to just the name.ext when this key/value will
	// be logged.
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	// Construct the slog JSON handler for use.
	handler := slog.Handler(
		slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       slog.Level(minLevel),
			ReplaceAttr: f,
		}))
	// If events are to be processed, wrap the JSON handler around the custom
	// log handler.
	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		handler = newLogHandler(handler, events)
	}

	// inject the sevice to the attrs
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}
	handler = handler.WithAttrs(attrs)

	return &Logger{
		handler:      handler,
		tracerIDFunc: tracerIDFunc,
	}

}
