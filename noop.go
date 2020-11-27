package ctxd

import "context"

// NoOpLogger is a contextualized logger stub.
type NoOpLogger struct{}

var _ Logger = NoOpLogger{}

// Debug discards debug message.
func (NoOpLogger) Debug(_ context.Context, _ string, _ ...interface{}) {}

// Info discards informational message.
func (NoOpLogger) Info(_ context.Context, _ string, _ ...interface{}) {}

// Important discards important message.
func (NoOpLogger) Important(_ context.Context, _ string, _ ...interface{}) {}

// Warn discards warning message.
func (NoOpLogger) Warn(_ context.Context, _ string, _ ...interface{}) {}

// Error discards error message.
func (NoOpLogger) Error(_ context.Context, _ string, _ ...interface{}) {}

// CtxdLogger is a provider.
func (NoOpLogger) CtxdLogger() Logger {
	return NoOpLogger{}
}
