package zaputil

import (
	"context"

	"go.uber.org/zap"
)

var (
	ContextLoggerKey        = "logger"
	ContextSugaredLoggerKey = "sugared-logger"
)

func PackLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}

func PackSugaredLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ContextSugaredLoggerKey, logger)
}

func UnpackLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLoggerKey).(*zap.Logger)
	if !ok {
		return nil
	}

	return logger
}

func UnpackSugaredLogger(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(ContextSugaredLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return nil
	}

	return logger
}
