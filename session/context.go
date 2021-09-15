package session

import "context"

var (
	ContextKey = "client"
)

func PackContext(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, ContextKey, client)
}

func UnpackContext(ctx context.Context) *Client {
	logger, ok := ctx.Value(ContextKey).(*Client)
	if !ok {
		return nil
	}

	return logger
}
