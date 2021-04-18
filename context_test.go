package ctxd_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
)

func TestSetFields(t *testing.T) {
	ctx := context.Background()

	ctx = ctxd.AddFields(ctx, "foo", 1, "bar", 1, "baz", 1)
	octx := ctx
	assert.Equal(t, []interface{}{"foo", 1, "bar", 1, "baz", 1}, ctxd.Fields(ctx))

	ctx = ctxd.SetFields(ctx, "quux", 2, "bar", 2, "foo", 2)
	assert.Equal(t, []interface{}{"foo", 2, "bar", 2, "baz", 1, "quux", 2}, ctxd.Fields(ctx))
	assert.Equal(t, []interface{}{"foo", 1, "bar", 1, "baz", 1}, ctxd.Fields(octx))

	ctx = ctxd.ClearFields(ctx)
	assert.Equal(t, []interface{}(nil), ctxd.Fields(ctx))
	assert.Equal(t, []interface{}{"foo", 1, "bar", 1, "baz", 1}, ctxd.Fields(octx))
}

func BenchmarkSetFields(b *testing.B) {
	ctx := ctxd.AddFields(context.Background(), "foo", 1, "bar", 1, "baz", 1)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ctxd.SetFields(ctx, "quux", 2, "bar", 2, "foo", 2)
	}
}
