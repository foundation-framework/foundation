package log

import (
	"context"
)

var (
	ContextLoggerKey = "logger"
)

func PackLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}

func UnpackLogger(ctx context.Context) *Logger {
	logger, ok := ctx.Value(ContextLoggerKey).(*Logger)
	if !ok {
		return nil
	}

	return logger
}
