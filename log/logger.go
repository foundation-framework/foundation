package log

import "context"

var (
	ContextKey = "logger"
)

type Logger interface {
	Debugf(template string, args interface{})
	Infof(template string, args interface{})
	Warnf(template string, args interface{})
	Errorf(template string, args interface{})
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
