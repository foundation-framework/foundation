package zapext

import (
	"context"

	"go.uber.org/zap"
)

const (
	DefaultContextLoggerKey = "logger"
)

func PackLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, DefaultContextLoggerKey, logger)
}

func UnpackLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(DefaultContextLoggerKey).(*zap.Logger)
	if !ok {
		return nil
	}

	return logger
}

func UnpackLoggerName(ctx context.Context, name string) *zap.Logger {
	logger, ok := ctx.Value(name).(*zap.Logger)
	if !ok {
		return nil
	}

	return logger
}

func UnpackLoggerNop(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(DefaultContextLoggerKey).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}

func UnpackLoggerNameNop(ctx context.Context, name string) *zap.Logger {
	logger, ok := ctx.Value(name).(*zap.Logger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}
