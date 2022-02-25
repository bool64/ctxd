# Contextualized Structured Logging and Error Handling for Go


This library provides context-driven structured error and logger contracts.

[![Build Status](https://github.com/bool64/ctxd/workflows/test/badge.svg)](https://github.com/bool64/ctxd/actions?query=branch%3Amaster+workflow%3Atest)
[![Coverage Status](https://codecov.io/gh/bool64/ctxd/branch/master/graph/badge.svg)](https://codecov.io/gh/bool64/ctxd)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/bool64/ctxd)
[![time tracker](https://wakatime.com/badge/github/bool64/ctxd.svg)](https://wakatime.com/badge/github/bool64/ctxd)
![Code lines](https://sloc.xyz/github/bool64/ctxd/?category=code)
![Comments](https://sloc.xyz/github/bool64/ctxd/?category=comments)

## Usage

* Create an adapter for your logger that implements `ctxd.Logger` or use [`zapctxd`](https://github.com/bool64/zapctxd)
that is built around awesome [`go.uber.org/zap`](https://pkg.go.dev/go.uber.org/zap).
* Add fields to context and pass it around.
* Use context for last-mile logging or error emitting.

## Example

### Structured Logging

```go
logger := ctxd.LoggerMock{}

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
```

Logger can be instrumented with persistent fields that are affecting every context.

```go
lm := ctxd.LoggerMock{}

var globalLogger ctxd.Logger = &lm

localLogger := ctxd.LoggerWithFields(globalLogger, "local", 123)

ctx1 := ctxd.AddFields(context.Background(),
    "ctx", 1,
    "foo", "bar",
)
ctx2 := ctxd.AddFields(context.Background(), "ctx", 2)

localLogger.Info(ctx1, "hello", "he", "lo")
localLogger.Warn(ctx2, "bye", "by", "ee")

fmt.Print(lm.String())

// Output:
// info: hello {"ctx":1,"foo":"bar","he":"lo","local":123}
// warn: bye {"by":"ee","ctx":2,"local":123}
```

### Handling Errors

```go
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
    lm                 = ctxd.LoggerMock{}
    logger ctxd.Logger = &lm
)

// Inspect error fields.
var se ctxd.StructuredError
if errors.As(err, &se) {
    fmt.Printf("error fields: %v\n", se.Fields())
}

// Log errors.
ctxd.LogError(ctx, err, logger.Error)
ctxd.LogError(ctx, err2, logger.Warn)
fmt.Print(lm.String())

// Output:
// error fields: map[field1:1 field2:abc field3:3]
// error: something failed {"field1":1,"field2":"abc","field3":3}
// warn: wrapped: something failed {"field1":1,"field2":"abc","field3":3,"field4":true,"field5":"V"}
```

