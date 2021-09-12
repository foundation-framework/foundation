package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var (
	ContextKey = "logger"
)

func init() {
	if err := zap.RegisterSink(udpScheme, newUdpSink); err != nil {
		panic(fmt.Errorf("unexpected error: %s", err.Error()))
	}
}

func PackContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ContextKey, logger)
}

func UnpackContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(ContextKey).(Logger)
	if !ok {
		return nil
	}

	return logger
}
