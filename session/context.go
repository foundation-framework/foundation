package session

import "context"

var (
	ContextClientKey = "client"
)

func PackClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, ContextClientKey, client)
}

func UnpackClient(ctx context.Context) *Client {
	logger, ok := ctx.Value(ContextClientKey).(*Client)
	if !ok {
		return nil
	}

	return logger
}
