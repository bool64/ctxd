package ctxd

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
)

// LogFunc defines contextualized logger function.
type LogFunc func(ctx context.Context, msg string, keysAndValues ...interface{})

// LogError pushes error value to a contextualized logger method.
//
// If err is nil, LogError produces no operation.
// LogError function matches Logger methods, e.g. Error.
func LogError(ctx context.Context, err error, l LogFunc) {
	if err == nil {
		return
	}

	var se StructuredError

	if errors.As(err, &se) {
		// Discarding keys and values from context as error already has full set of fields prepared on invocation.
		l(ClearFields(ctx), se.Error(), se.Tuples()...)

		return
	}

	l(ctx, err.Error())
}

// StructuredError defines error with message and data.
type StructuredError interface {
	// Error returns message of error.
	Error() string

	// Tuples returns structured data of error in form of loosely-typed key-value pairs.
	Tuples() []interface{}

	// Fields returns structured data of error as a map.
	Fields() map[string]interface{}
}

type wrappedError struct {
	message string
	err     error
}

func (we wrappedError) Unwrap() error {
	return we.err
}

func (we wrappedError) Error() string {
	return we.message + ": " + we.err.Error()
}

// WrapError returns an error annotated with message and structured data.
//
// If err is nil, WrapError returns nil.
// LogError fields from context are also added to error structured data.
func WrapError(ctx context.Context, err error, message string, keysAndValues ...interface{}) error {
	if err == nil {
		return nil
	}

	if message != "" {
		err = wrappedError{
			err:     err,
			message: message,
		}
	}

	se, ok := newError(ctx, err, keysAndValues...)
	if ok {
		return wrappedStructuredError{
			structuredError: se,
		}
	}

	return err
}

// NewError creates error with optional structured data.
//
// LogError fields from context are also added to error structured data.
func NewError(ctx context.Context, message string, keysAndValues ...interface{}) error {
	// nolint:goerr113 // Static errors can be used with WrapError.
	err := errors.New(message)

	se, ok := newError(ctx, err, keysAndValues...)
	if ok {
		return se
	}

	return err
}

// Tuples is a slice of keys and values, e.g. {"key1", 1, "key2", "val2"}.
type Tuples []interface{}

type structuredError struct {
	err           error
	keysAndValues Tuples
}

type wrappedStructuredError struct {
	structuredError
}

// Unwrap implements errors wrapper.
func (wse wrappedStructuredError) Unwrap() error {
	return wse.err
}

// Fields creates a map from key-value pairs.
func (t Tuples) Fields() map[string]interface{} {
	result := make(map[string]interface{}, len(t))

	var (
		label string
		ok    bool
	)

	for i, l := range t {
		if label == "" {
			label, ok = l.(string)
			if !ok {
				result["malformedFields"] = []interface{}(t[i:])

				break
			}
		} else {
			result[label] = l
			label = ""
		}
	}

	return result
}

// Fields returns structured data of error as a map.
func (se structuredError) Fields() map[string]interface{} {
	return se.keysAndValues.Fields()
}

// Error returns message and data serialized to a string.
func (se structuredError) String() string {
	err := se.err.Error()

	var (
		label string
		ok    bool
	)

	for i, l := range se.keysAndValues {
		if label == "" {
			label, ok = l.(string)
			if !ok {
				err += fmt.Sprintf(", malformed fields: %+v", se.keysAndValues[i:])

				break
			}
		} else {
			err += fmt.Sprintf(", %s: %+v", label, l)
			label = ""
		}
	}

	return err
}

// Error returns message of error.
func (se structuredError) Error() string {
	return se.err.Error()
}

// KeysAndValues returns structured data of error in form of loosely-typed key-value pairs.
func (se structuredError) Tuples() []interface{} {
	return se.keysAndValues[0:len(se.keysAndValues):len(se.keysAndValues)]
}

func newError(ctx context.Context, err error, keysAndValues ...interface{}) (structuredError, bool) {
	var (
		se        StructuredError
		kv        = keysAndValues
		tuples    []interface{}
		ctxFields []interface{}
	)

	if errors.As(err, &se) {
		tuples = se.Tuples()
	}

	ctxFields = Fields(ctx)

	if len(tuples)+len(ctxFields) > 0 {
		kv = make([]interface{}, 0, len(kv)+len(tuples)+len(ctxFields))

		kv = append(kv, tuples...)
		kv = append(kv, keysAndValues...)
		kv = append(kv, ctxFields...)
	}

	if len(kv) > 1 {
		return structuredError{
			err:           err,
			keysAndValues: kv,
		}, true
	}

	return structuredError{}, false
}

var (
	_ encoding.TextMarshaler = structuredError{}
	_ json.Marshaler         = structuredError{}
)

func (se structuredError) MarshalText() ([]byte, error) {
	return []byte(se.err.Error()), nil
}

func (se structuredError) MarshalJSON() ([]byte, error) {
	return json.Marshal(se.err.Error())
}
