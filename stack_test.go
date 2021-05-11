package ctxd_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStack_StackTrace(t *testing.T) {
	err := ctxd.NewError(context.Background(), "failed")

	var e interface {
		StackTrace() ctxd.StackTrace
	}

	assert.True(t, errors.As(err, &e))
	require.NotNil(t, e)
	assert.Equal(t, "stack_test.go:15", fmt.Sprintf("%v", e.StackTrace()[0]))
}
