package ctxd_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMock_Error(t *testing.T) {
	var e error

	m := ctxd.LoggerMock{
		OnError: func(err error) {
			e = err
		},
	}

	m.Error(context.Background(), "foo", "func", func() {})
	assert.EqualError(t, e, "json: unsupported type: func()")
	assert.Equal(t, "", m.String())

	m.OnError = nil

	m.Error(context.Background(), "foo", "func", func() {})
	assert.Equal(t, "json: unsupported type: func()\n", m.String())
}

func TestLoggerMock_Debug(t *testing.T) {
	m := ctxd.LoggerMock{}

	ctx := ctxd.AddFields(context.Background(), "foo", 1, "bar", 2)

	m.Debug(ctx, "debug message", "baz", 3)
	m.Info(ctx, "info message", "baz", 3)
	m.Warn(ctx, "warn message", "baz", 3)
	m.Error(ctx, "error message", "baz", 3)
	m.Important(ctx, "important message", "baz", 3)

	assert.Equal(t, `debug: debug message {"bar":2,"baz":3,"foo":1}
info: info message {"bar":2,"baz":3,"foo":1}
warn: warn message {"bar":2,"baz":3,"foo":1}
error: error message {"bar":2,"baz":3,"foo":1}
important: important message {"bar":2,"baz":3,"foo":1}`+"\n", m.String(), m.String())

	data := map[string]interface{}{"bar": 2, "baz": 3, "foo": 1}

	assert.Equal(t, "debug message", m.LoggedEntries[0].Message)
	assert.Equal(t, data, m.LoggedEntries[0].Data)
	assert.Equal(t, "info message", m.LoggedEntries[1].Message)
	assert.Equal(t, data, m.LoggedEntries[1].Data)
	assert.Equal(t, "warn message", m.LoggedEntries[2].Message)
	assert.Equal(t, data, m.LoggedEntries[2].Data)
	assert.Equal(t, "error message", m.LoggedEntries[3].Message)
	assert.Equal(t, data, m.LoggedEntries[3].Data)
	assert.Equal(t, "important message", m.LoggedEntries[4].Message)
	assert.Equal(t, data, m.LoggedEntries[4].Data)
}
