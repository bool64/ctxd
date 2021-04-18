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

// SetFields returns context with added loosely-typed key-value pairs as fields.
//
// Values of same keys that are already existing in parent context are replaced.
func SetFields(ctx context.Context, keysAndValues ...interface{}) context.Context {
	fields := Fields(ctx)

	var existing []int

	for i := 0; i < len(keysAndValues); i += 2 {
		for j := 0; j < len(fields); j += 2 {
			if keysAndValues[i] == fields[j] {
				existing = append(existing, j)
			}
		}
	}

	if len(existing) == 0 {
		return context.WithValue(ctx, fieldsCtxKey{}, append(fields, keysAndValues...))
	}

	// Create result slice that will fit a copy of fields and will have a capacity for new keys and values.
	result := make([]interface{}, len(fields), len(fields)+len(keysAndValues)-2*len(existing))
	copy(result, fields)

	for i := 0; i < len(keysAndValues); i += 2 {
		found := false

		for _, e := range existing {
			if result[e] == keysAndValues[i] {
				result[e+1] = keysAndValues[i+1]
				found = true

				break
			}
		}

		if found {
			continue
		}

		result = append(result, keysAndValues[i], keysAndValues[i+1])
	}

	return context.WithValue(ctx, fieldsCtxKey{}, result)
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
