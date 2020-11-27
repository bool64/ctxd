package ctxd_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
)

func TestNoOpLogger_CtxdLogger(_ *testing.T) {
	n := ctxd.NoOpLogger{}
	ctx := context.Background()

	n.Debug(ctx, "msg")
	n.Info(ctx, "msg")
	n.Warn(ctx, "msg")
	n.Error(ctx, "msg")
	n.Important(ctx, "msg")
	n.CtxdLogger()
}
