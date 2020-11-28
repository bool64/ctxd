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
