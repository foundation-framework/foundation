package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var (
	ContextKey = "logger"
)

type Logger interface {
	Named(name string) Logger
	With(keysAndValues ...interface{}) Logger

	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	DPanic(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

func init() {
	fmt.Println("Initializing logger package")

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
