package ctxd

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
)

// LoggerMock logs messages to internal buffer.
type LoggerMock struct {
	mu sync.Mutex
	bytes.Buffer
	OnError func(err error)
}

func (m *LoggerMock) failed(err error) bool {
	if err == nil {
		return false
	}

	if m.OnError != nil {
		m.OnError(err)

		return true
	}

	_, _ = m.WriteString(err.Error() + "\n")

	return true
}

func (m *LoggerMock) log(ctx context.Context, level, msg string, keysAndValues []interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	jm, err := json.Marshal(Tuples(append(Fields(ctx), keysAndValues...)).Fields())
	if m.failed(err) {
		return
	}

	out := LogWriter(ctx)
	if out == nil {
		out = m
	}

	if IsDebug(ctx) {
		_, err = out.Write([]byte("debug mode, "))
		if m.failed(err) {
			return
		}
	}

	_, err = out.Write([]byte(level + ": " + msg + " "))
	if m.failed(err) {
		return
	}

	_, err = out.Write(jm)
	if m.failed(err) {
		return
	}

	_, err = out.Write([]byte("\n"))
	if m.failed(err) {
		return
	}
}

// Debug logs a message.
func (m *LoggerMock) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.log(ctx, "debug", msg, keysAndValues)
}

// Info logs a message.
func (m *LoggerMock) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.log(ctx, "info", msg, keysAndValues)
}

// Important logs a message.
func (m *LoggerMock) Important(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.log(ctx, "important", msg, keysAndValues)
}

// Warn logs a message.
func (m *LoggerMock) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.log(ctx, "warn", msg, keysAndValues)
}

// Error logs a message.
func (m *LoggerMock) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.log(ctx, "error", msg, keysAndValues)
}
