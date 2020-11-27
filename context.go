package ctxd

import (
	"context"
)

type fieldsCtxKey struct{}

// AddFields returns context with added loosely-typed key-value pairs as fields.
//
// If key-value pairs exist in parent context already, new pairs are appended.
func AddFields(ctx context.Context, keysAndValues ...interface{}) context.Context {
	return context.WithValue(ctx, fieldsCtxKey{}, append(Fields(ctx), keysAndValues...))
}

// ClearFields returns context without any fields.
func ClearFields(ctx context.Context) context.Context {
	_, ok := ctx.Value(fieldsCtxKey{}).([]interface{})
	if !ok {
		return ctx
	}

	return context.WithValue(ctx, fieldsCtxKey{}, nil)
}

// Fields returns loosely-typed key-value pairs found in context or nil.
func Fields(ctx context.Context) []interface{} {
	keysAndValues, ok := ctx.Value(fieldsCtxKey{}).([]interface{})
	if !ok {
		return nil
	}

	return keysAndValues[:len(keysAndValues):len(keysAndValues)]
}
