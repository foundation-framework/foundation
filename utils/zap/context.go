package zaputil

import (
	"context"

	"go.uber.org/zap"
)

var (
	ContextLoggerKey = "logger"
)

func PackLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}

func UnpackLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLoggerKey).(*zap.Logger)
	if !ok {
		return nil
	}

	return logger
}

func UnpackLoggerNop(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLoggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}
