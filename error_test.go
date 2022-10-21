package ctxd_test

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/usecase/status"
)

func TestWrap(t *testing.T) {
	ctx := context.Background()
	ctx = ctxd.AddFields(ctx, "country", "us")

	var (
		stringer fmt.Stringer
		err      = ctxd.WrapError(ctx, status.NotFound, "failed to find order", "id", 123)
	)

	assert.NotNil(t, err)
	assert.Equal(t, "failed to find order: not found", err.Error())
	assert.True(t, errors.As(err, &stringer))
	assert.Equal(t, "failed to find order: not found, id: 123, country: us", stringer.String())

	logOut := &bytes.Buffer{}
	ctx = ctxd.WithLogWriter(ctx, logOut)

	logger := ctxd.LoggerMock{}

	ctxd.LogError(ctx, err, logger.Error)

	var st status.Code

	assert.True(t, errors.As(err, &st))
	assert.Equal(t, st, status.NotFound)

	assert.True(t, errors.Is(err, status.NotFound))
	assert.False(t, errors.Is(err, status.Unknown))

	assert.Equal(t,
		`error: failed to find order: not found {"country":"us","id":123}`+"\n",
		logOut.String())
}

func TestWrap_noKeys(t *testing.T) {
	err := ctxd.WrapError(context.Background(), errors.New("failed"), "unable to can")
	assert.NotNil(t, err)
	assert.Equal(t, "unable to can: failed", err.Error())
}

func TestWrap_nilErr(t *testing.T) {
	err := ctxd.WrapError(context.Background(), nil, "failed to win")
	assert.Nil(t, err)
}

func TestWrap_noCtxKeys(t *testing.T) {
	var (
		stringer fmt.Stringer
		err      = ctxd.WrapError(context.Background(), errors.New("failed"), "unable to can",
			"key1", 123,
			"key2", "abc",
		)
	)

	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &stringer))
	assert.Equal(t, "unable to can: failed, key1: 123, key2: abc", stringer.String())
}

func TestWrap_doubleWrap(t *testing.T) {
	ctx := context.Background()
	err := ctxd.WrapError(ctx, status.NotFound, "failed to find order", "id", 123)

	ctxd.LogError(ctx, err, func(ctx context.Context, msg string, keysAndValues ...interface{}) {
		assert.Equal(t, "failed to find order: not found", msg)
		assert.Equal(t, []interface{}{"id", 123}, keysAndValues)
	})

	ctx = ctxd.AddFields(ctx, "extra", 321)
	err = ctxd.WrapError(ctx, err, "wrapped")
	ctxd.LogError(ctx, err, func(ctx context.Context, msg string, keysAndValues ...interface{}) {
		assert.Equal(t, "wrapped: failed to find order: not found", msg)
		assert.Equal(t, []interface{}{"id", 123, "extra", 321}, keysAndValues)
	})
}

func TestNew(t *testing.T) {
	var (
		stringer fmt.Stringer
		err      = ctxd.NewError(context.Background(), "failed",
			"key1", 123,
			"key2", "abc",
		)
	)

	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &stringer))
	assert.Equal(t, "failed, key1: 123, key2: abc", stringer.String())
}

func TestNew_noFields(t *testing.T) {
	err := ctxd.NewError(context.Background(), "failed")
	assert.NotNil(t, err)
	assert.Equal(t, "failed", err.Error())
}

func TestNew_malformedFields(t *testing.T) {
	var (
		stringer fmt.Stringer
		err      = ctxd.NewError(context.Background(), "failed",
			"key1", 1,
			123, 2, // non-string key is breaking processing
			"key3", 3,
		)
	)

	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &stringer))
	assert.Equal(t, "failed, key1: 1, malformed fields: [123 2 key3 3]", stringer.String())
}

func TestStructuredError_Fields(t *testing.T) {
	var (
		se  ctxd.StructuredError
		err = ctxd.NewError(context.Background(), "failed",
			"key1", 1,
			123, 2, // Non-string key is breaking processing.
			"key3", 3,
		)
	)

	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &se))
	assert.Equal(t, map[string]interface{}{
		"key1": 1,
		"malformedFields": []interface{}{
			123, 2, // Non-string key is breaking processing.
			"key3", 3,
		},
	}, se.Fields())

	j, jerr := json.Marshal(err)
	require.NoError(t, jerr)
	assert.Equal(t, `"failed"`, string(j))
	assert.Equal(t, `failed`, fmt.Sprintf("%s", err))

	var tm encoding.TextMarshaler

	assert.True(t, errors.As(err, &tm))

	m, merr := tm.MarshalText()

	require.NoError(t, merr)
	assert.Equal(t, `failed`, string(m))
}

func TestLog(t *testing.T) {
	ctx := context.Background()
	logOut := &bytes.Buffer{}
	ctx = ctxd.WithLogWriter(ctx, logOut)

	logger := ctxd.LoggerMock{}

	err := ctxd.NewError(ctx, "failed", "id", 123)
	assert.NotNil(t, err)

	ctxd.LogError(ctx, err, logger.Error)
	ctxd.LogError(ctx, nil, logger.Error)
	ctxd.LogError(ctx, ctxd.NewError(ctx, "failed with no fields"), logger.Error)
	assert.Equal(t,
		`error: failed {"id":123}
error: failed with no fields null
`,
		logOut.String())
}

func BenchmarkStructuredError_Error(b *testing.B) {
	ctx := context.Background()
	ctx = ctxd.AddFields(ctx, "country", "us")

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := ctxd.WrapError(ctx, status.NotFound, "failed to find order", "id", 123)
		if err != nil {
			_ = err.Error()
		} else {
			b.Fail()
		}
	}
}

func BenchmarkFmt_Errorf(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := fmt.Errorf("failed to find item, country: %s, id: %d: %w", "us", 123, status.NotFound)
		_ = err.Error()
	}
}

func BenchmarkWrapError(b *testing.B) {
	ctx := context.Background()
	ctx = ctxd.AddFields(ctx, "country", "us")

	logFunc := func(ctx context.Context, msg string, keysAndValues ...interface{}) {}

	e1 := ctxd.NewError(ctx, "not found", "a", 1)

	b.ReportAllocs()

	var err error
	for i := 0; i < b.N; i++ {
		err = ctxd.WrapError(ctx, e1, "failed to find order", "id", 123)
	}
	ctxd.LogError(ctx, err, logFunc)
}

func TestSentinelError_Error(t *testing.T) {
	assert.EqualError(t, ctxd.SentinelError("failed"), "failed")
}

func TestLabeledError(t *testing.T) {
	err1 := errors.New("failed")
	label1 := ctxd.SentinelError("miserably")
	label2 := ctxd.SentinelError("hopelessly")

	err := ctxd.LabeledError(fmt.Errorf("oops: %w", err1), label1, label2)

	assert.True(t, errors.Is(err, err1))
	assert.True(t, errors.Is(err, label1))
	assert.True(t, errors.Is(err, label2))

	// Labels do not implicitly contribute to error message.
	assert.Equal(t, "oops: failed", err.Error())

	// If there are two matches, only first is returned.
	var se ctxd.SentinelError

	assert.True(t, errors.As(err, &se))
	assert.Equal(t, "miserably", string(se))
}

func TestTuples_Fields(t *testing.T) {
	assert.Equal(t, map[string]interface{}(nil), ctxd.Tuples(nil).Fields()) // Empty tuples.
	assert.Equal(t, map[string]interface{}{"malformedFields": []interface{}{1, 2}},
		ctxd.Tuples{1, 2}.Fields()) // String key expected.
	assert.Equal(t, map[string]interface{}{"malformedFields": []interface{}{"key"}},
		ctxd.Tuples{"key"}.Fields()) // Key without a value.
	assert.Equal(t, map[string]interface{}{"malformedFields": []interface{}{"", 123}},
		ctxd.Tuples{"", 123}.Fields()) // Empty key.
	assert.Equal(t, map[string]interface{}{"a": 123, "b": 456},
		ctxd.Tuples{"a", 123, "b", 456}.Fields()) // All good.
}

func TestNewMulti(t *testing.T) {
	errPrimary := errors.New("failed")
	errSecondary1 := ctxd.SentinelError("miserably")
	errSecondary2 := ctxd.SentinelError("hopelessly")
	errSecondary3 := ctxd.SentinelError("deadly")

	err := ctxd.MultiError(fmt.Errorf("oops: %w", errPrimary), errSecondary1, errSecondary2)

	assert.True(t, errors.Is(err, errPrimary))
	assert.True(t, errors.Is(err, errSecondary1))
	assert.True(t, errors.Is(err, errSecondary2))
	assert.False(t, errors.Is(err, errSecondary3))

	// Labels do not implicitly contribute to error message.
	assert.Equal(t, "oops: failed", err.Error())

	// If there are two matches, only first is returned.
	var errSentinel ctxd.SentinelError

	assert.True(t, errors.As(err, &errSentinel))
	assert.Equal(t, "miserably", string(errSentinel))
}
