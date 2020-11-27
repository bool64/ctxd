package ctxd_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/bool64/ctxd"
)

func ExampleAddFields() {
	logger := testLogger{}

	// Once instrumented context can aid logger with structured information.
	ctx := ctxd.AddFields(context.Background(), "foo", "bar")

	logger.Info(ctx, "something happened")

	// Also context contributes additional information to structured errors.
	err := ctxd.NewError(ctx, "something failed", "baz", "quux")

	ctxd.LogError(ctx, err, logger.Error)

	fmt.Print(logger.String())
	// Output:
	// info: something happened {"foo":"bar"}
	// error: something failed {"baz":"quux","foo":"bar"}
}

func ExampleLoggerWithFields() {
	tl := testLogger{}

	var globalLogger ctxd.Logger = &tl

	localLogger := ctxd.LoggerWithFields(globalLogger, "local", 123)

	ctx1 := ctxd.AddFields(context.Background(),
		"ctx", 1,
		"foo", "bar",
	)
	ctx2 := ctxd.AddFields(context.Background(), "ctx", 2)

	localLogger.Info(ctx1, "hello", "he", "lo")
	localLogger.Warn(ctx2, "bye", "by", "ee")

	fmt.Print(tl.String())

	// Output:
	// info: hello {"ctx":1,"foo":"bar","he":"lo","local":123}
	// warn: bye {"by":"ee","ctx":2,"local":123}
}

func ExampleWrapError() {
	ctx := context.Background()

	// Elaborate context with fields.
	ctx = ctxd.AddFields(ctx,
		"field1", 1,
		"field2", "abc",
	)

	// Add more fields when creating error.
	err := ctxd.NewError(ctx, "something failed",
		"field3", 3.0,
	)

	err2 := ctxd.WrapError(
		// You can use same or any other context when wrapping error.
		ctxd.AddFields(context.Background(), "field5", "V"),
		err, "wrapped",
		"field4", true)

	// Setup your logger.
	var (
		tl                 = testLogger{}
		logger ctxd.Logger = &tl
	)

	// Inspect error fields.
	var se ctxd.StructuredError
	if errors.As(err, &se) {
		fmt.Printf("error fields: %v\n", se.Fields())
	}

	// Log errors.
	ctxd.LogError(ctx, err, logger.Error)
	ctxd.LogError(ctx, err2, logger.Warn)
	fmt.Print(tl.String())

	// Output:
	// error fields: map[field1:1 field2:abc field3:3]
	// error: something failed {"field1":1,"field2":"abc","field3":3}
	// warn: wrapped: something failed {"field1":1,"field2":"abc","field3":3,"field4":true,"field5":"V"}
}
