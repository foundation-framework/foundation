package validatorutil

import (
	"context"

	"github.com/intale-llc/foundation/net/sockets"
)

func SocketsHandler(fn func(err error) (string, interface{})) sockets.HandlerFunc {
	return func(ctx context.Context, data interface{}) (string, interface{}) {
		if err := Struct(data); err != nil {
			return fn(err)
		}

		return "", nil
	}
}
