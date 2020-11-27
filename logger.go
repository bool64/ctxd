package ctxd

import (
	"context"
	"io"
	"sync"
)

// Logger is a contextualized structured logger.
//
// Logging methods accept keys and values as variadic argument that contains loosely-typed key-value pairs.
// When processing pairs, the first element of the pair is used as the field key and the second as the field value.
type Logger interface {
	// Debug logs a message.
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})

	// Info logs a message.
	Info(ctx context.Context, msg string, keysAndValues ...interface{})

	// Important forcibly logs an important message with level INFO disregarding logger level constraints.
	// Can be used for logging historically important information.
	Important(ctx context.Context, msg string, keysAndValues ...interface{})

	// Warn logs a message.
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})

	// Error logs a message.
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

// LoggerProvider is an embeddable provider interface.
type LoggerProvider interface {
	CtxdLogger() Logger
}

type (
	isDebugCtxKey   struct{}
	logWriterCtxKey struct{}
)

type syncWriter struct {
	m sync.Mutex
	w io.Writer
}

func (sw *syncWriter) Write(p []byte) (n int, err error) {
	sw.m.Lock()
	defer sw.m.Unlock()

	return sw.w.Write(p)
}

// WithLogWriter returns context with custom log writer.
// Can be useful to write logs into response stream.
func WithLogWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, logWriterCtxKey{}, &syncWriter{w: w})
}

// LogWriter returns custom log writer found in context or nil.
func LogWriter(ctx context.Context) io.Writer {
	w, ok := ctx.Value(logWriterCtxKey{}).(*syncWriter)
	if !ok {
		return nil
	}

	return w
}

// WithDebug returns context with debug flag enabled.
func WithDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, isDebugCtxKey{}, true)
}

// IsDebug returns true if debug flag is enabled in context.
func IsDebug(ctx context.Context) bool {
	_, ok := ctx.Value(isDebugCtxKey{}).(bool)

	return ok
}

// LoggerWithFields instruments contextualized logger with global fields.
func LoggerWithFields(logger Logger, keysAndValues ...interface{}) Logger {
	return &withFields{
		logger:        logger,
		keysAndValues: keysAndValues,
	}
}

type withFields struct {
	logger        Logger
	keysAndValues []interface{}
}

func (w *withFields) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	w.logger.Debug(AddFields(ctx, w.keysAndValues...), msg, keysAndValues...)
}

func (w *withFields) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	w.logger.Info(AddFields(ctx, w.keysAndValues...), msg, keysAndValues...)
}

func (w *withFields) Important(ctx context.Context, msg string, keysAndValues ...interface{}) {
	w.logger.Important(AddFields(ctx, w.keysAndValues...), msg, keysAndValues...)
}

func (w *withFields) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	w.logger.Warn(AddFields(ctx, w.keysAndValues...), msg, keysAndValues...)
}

func (w *withFields) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	w.logger.Error(AddFields(ctx, w.keysAndValues...), msg, keysAndValues...)
}
