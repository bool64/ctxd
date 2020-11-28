package ctxd_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
)

func TestWithFields(t *testing.T) {
	l := ctxd.LoggerMock{}
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
