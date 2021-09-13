package transport

import "context"

// Server describes server for real-time communications
type Server interface {

	// Listen starts listening for connections on provided address or adds handler
	// for a specified path (endpoint parameter implementation defined)
	Listen(endpoint string) error

	// Shutdown gracefully shutdowns the listener
	// All active connections will be closed
	Shutdown(ctx context.Context) error

	// OnConnection sets callback for an incoming connection
	// Multiple callbacks allowed
	OnConnection(func(Connection))

	// OnError sets a callback for non-critical connection errors
	// Multiple callbacks allowed
	OnError(func(error))

	// OnClose sets a callback for listener closure
	// Error can be nil
	//
	// Multiple callbacks allowed
	OnClose(func(error))
}
