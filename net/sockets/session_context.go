package sockets

import (
	"context"
)

const (
	SessionContextKey = "session"
)

func PackSession(ctx context.Context, client *Session) context.Context {
	return context.WithValue(ctx, SessionContextKey, client)
}

func UnpackSession(ctx context.Context) *Session {
	logger, ok := ctx.Value(SessionContextKey).(*Session)
	if !ok {
		return nil
	}

	return logger
}

func UnpackSessionNamed(ctx context.Context, name string) *Session {
	logger, ok := ctx.Value(name).(*Session)
	if !ok {
		return nil
	}

	return logger
}
