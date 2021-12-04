package sockets

import (
	"context"
)

const (
	SessionContextKey = "session"
)

func PackSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionContextKey, session)
}

func UnpackSession(ctx context.Context) *Session {
	session, ok := ctx.Value(SessionContextKey).(*Session)
	if !ok {
		return nil
	}

	return session
}

func UnpackSessionNamed(ctx context.Context, name string) *Session {
	session, ok := ctx.Value(name).(*Session)
	if !ok {
		return nil
	}

	return session
}
