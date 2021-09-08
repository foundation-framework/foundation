package rest

import (
	"context"

	"github.com/gorilla/mux"
)

// Listener describes basic REST server
type Listener interface {

	// Listen starts listening for requests on a provided address or path
	// (endpoint parameter implementation defined)
	Listen(endpoint string) error

	// Shutdown gracefully shutdowns server
	// All active connections will be closed
	Shutdown(ctx context.Context) error

	// Router returns server router
	// You can add any handlers via this router BEFORE server started listen
	//
	// Any changes made after server has started listening have no effect
	Router() *mux.Router

	// OnError sets callback for non-critical server errors
	// Multiple callbacks allowed
	OnError(func(error))

	// OnClose sets callback for server closure
	// Multiple callbacks allowed
	OnClose(func(error))
}
