package session

import (
	"context"
)

const (
	DefaultContextClientKey = "client"
)

func PackClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, DefaultContextClientKey, client)
}

func UnpackClient(ctx context.Context) *Client {
	logger, ok := ctx.Value(DefaultContextClientKey).(*Client)
	if !ok {
		return nil
	}

	return logger
}

func UnpackClientNamed(ctx context.Context, name string) *Client {
	logger, ok := ctx.Value(name).(*Client)
	if !ok {
		return nil
	}

	return logger
}
