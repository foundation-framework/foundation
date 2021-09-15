package log

import (
	"context"
)

var (
	ContextKey = "logger"
)

func PackContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ContextKey, logger)
}

func UnpackContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(ContextKey).(*Logger)
	if !ok {
		return nil
	}

	return logger
}
