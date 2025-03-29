package logging

import (
	"context"
)

// LogParams is the struct that contains the parameters to create a logger
type LogParams struct {
	// AppName is the name of your application, that will be used as a field in the log
	AppName string
	// DebugLevel is the level of the log, if true, the log will be in debug level
	DebugLevel bool
	// AddAttributesFromContext is a function that will be called to add attributes to the log.
	// it should return a key and a value, example: []any{"account_id", ctx.Value(AccountUUIDKey)}
	// example: when you call logger.Info(ctx, "message"), the logger will add the attributes returned by the function
	// it will look like this: {"time":"2020-01-01T00:00:00","level":"INFO","file":"main.go:10","msg":"main: message","account_id":"1234567890"}
	AddAttributesFromContext func(ctx context.Context) []LogField
	// LogToFile is a flag to indicate if the log should be printed to file or not
	LogToFile bool
}

// LogParams is the struct that contains the parameters to create a logger
func New(params LogParams) Logger {
	return newLogger(params)
}

// Logger is a wrapper of the zerolog library adding some extra functionality
type Logger interface {
	// Info logs a message with INFO level
	Info(ctx context.Context, msg string)
	// Infof logs a message with INFO level and format the message, like fmt.Printf
	Infof(ctx context.Context, msg string, args ...any)
	// Infow logs a message with INFO level and add key and values to the log, example: logger.Infow(ctx, "message", "key", "value")
	Infow(ctx context.Context, msg string, fields ...LogField)
	// Debug logs a message with DEBUG level
	Debug(ctx context.Context, msg string)
	// Debugf logs a message with DEBUG level and format the message, like fmt.Printf
	Debugf(ctx context.Context, msg string, args ...any)
	// Debugw logs a message with DEBUG level and add key and values to the log, example: logger.Debugw(ctx, "message", "key", "value")
	Debugw(ctx context.Context, msg string, fields ...LogField)
	// Warn logs a message with WARN level
	Warn(ctx context.Context, msg string)
	// Warnf logs a message with WARN level and format the message, like fmt.Printf
	Warnf(ctx context.Context, msg string, args ...any)
	// Warnw logs a message with WARN level and add key and values to the log, example: logger.Warnw(ctx, "message", "key", "value")
	Warnw(ctx context.Context, msg string, fields ...LogField)
	// Error logs a message with ERROR level
	Error(ctx context.Context, msg string)
	// Errorf logs a message with ERROR level and format the message, like fmt.Printf
	Errorf(ctx context.Context, msg string, args ...any)
	// Errorw logs a message with ERROR level and add key and values to the log, example: logger.Errorw(ctx, "message", "key", "value")
	Errorw(ctx context.Context, msg string, fields ...LogField)
	// Fatal logs a message with FATAL level and exit the application
	Fatal(ctx context.Context, msg string)
	// Fatalf logs a message with FATAL level and format the message, like fmt.Printf and exit the application
	Fatalf(ctx context.Context, msg string, args ...any)
	// Fatalw logs a message with FATAL level and add key and values to the log, example: logger.Fatalw(ctx, "message", "key", "value") and exit the application
	Fatalw(ctx context.Context, msg string, fields ...LogField)
	// Critical logs a message with CRITICAL level
	Critical(ctx context.Context, msg string)
	// Criticalf logs a message with CRITICAL level and format the message, like fmt.Printf
	Criticalf(ctx context.Context, msg string, args ...any)
	// Criticalw logs a message with CRITICAL level and add key and values to the log, example: logger.Criticalw(ctx, "message", "key", "value")
	Criticalw(ctx context.Context, msg string, fields ...LogField)
	// Print implements the SetLog function on mysql library
	Print(args ...any)
	// Printf implements the SetLog function on elasticsearch library
	Printf(msg string, v ...any)
}
