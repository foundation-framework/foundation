package rest

import (
	"context"

	"github.com/gorilla/mux"
)

type Listener interface {
	Listen(addr string) error
	Close(ctx context.Context) error

	Router() *mux.Router

	OnError(func(error))
	OnClose(func(error))
	CloseChan() <-chan error
}
