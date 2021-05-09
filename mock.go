package ctxd

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"time"
)

// LoggerMock logs messages to internal buffer.
type LoggerMock struct {
	OnError func(err error)

	sync.Mutex
	bytes.Buffer
	LoggedEntries []struct {
		Time    time.Time              `json:"time"`
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}
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
	m.Lock()
	defer m.Unlock()

	data := Tuples(append(Fields(ctx), keysAndValues...)).Fields()

	jm, err := json.Marshal(data)
	if m.failed(err) {
		return
	}

	m.LoggedEntries = append(m.LoggedEntries, struct {
		Time    time.Time              `json:"time"`
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data,omitempty"`
	}{Time: time.Now(), Level: level, Message: msg, Data: data})

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
