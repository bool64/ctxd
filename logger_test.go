package ctxd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testLogger struct {
	bytes.Buffer
	TB testing.TB
}

func (t *testLogger) log(ctx context.Context, level, msg string, keysAndValues []interface{}) {
	jm, err := json.Marshal(ctxd.Tuples(append(ctxd.Fields(ctx), keysAndValues...)).Fields())
	require.NoError(t.TB, err)

	out := ctxd.LogWriter(ctx)
	if out == nil {
		out = t
	}

	if ctxd.IsDebug(ctx) {
		_, err = out.Write([]byte("debug mode, "))
		require.NoError(t.TB, err)
	}

	_, err = out.Write([]byte(level + ": " + msg + " "))
	require.NoError(t.TB, err)
	_, err = out.Write(jm)
	require.NoError(t.TB, err)
	_, err = out.Write([]byte("\n"))
	require.NoError(t.TB, err)
}

func (t *testLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	t.log(ctx, "debug", msg, keysAndValues)
}

func (t *testLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	t.log(ctx, "info", msg, keysAndValues)
}

func (t *testLogger) Important(ctx context.Context, msg string, keysAndValues ...interface{}) {
	t.log(ctx, "important", msg, keysAndValues)
}

func (t *testLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	t.log(ctx, "warn", msg, keysAndValues)
}

func (t *testLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	t.log(ctx, "error", msg, keysAndValues)
}

func TestWithFields(t *testing.T) {
	l := testLogger{}
	w := ctxd.LoggerWithFields(&l, "key1", 1, "key2", "abc")
	ctx := ctxd.AddFields(context.Background(), "key3", 3)

	w.Debug(ctx, "debug", "key4", 4)
	w.Info(ctx, "info", "key4", 4)
	w.Important(ctx, "important", "key4", 4)
	w.Warn(ctx, "warn", "key4", 4)
	w.Error(ctx, "error", "key4", 4)
	w.Info(ctxd.WithDebug(ctx), "info with debug", "key4", 4)

	assert.Equal(t, `debug: debug {"key1":1,"key2":"abc","key3":3,"key4":4}
info: info {"key1":1,"key2":"abc","key3":3,"key4":4}
important: important {"key1":1,"key2":"abc","key3":3,"key4":4}
warn: warn {"key1":1,"key2":"abc","key3":3,"key4":4}
error: error {"key1":1,"key2":"abc","key3":3,"key4":4}
debug mode, info: info with debug {"key1":1,"key2":"abc","key3":3,"key4":4}
`, l.String())
}
