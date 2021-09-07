package sockets

import "context"

// Listener describes a connection listener for real-time communications
type Listener interface {

	// Listen starts listening on specified address or adds handler for specified path
	// Implementation defined
	Listen(endpoint string) error

	// Close gracefully shutdowns the listener.
	// IMPORTANT!!! All active connections will be closed
	Close(ctx context.Context) error

	// OnConnection sets callback for an incoming connection
	// Multiple callbacks allowed
	OnConnection(func(Connection))

	// OnError sets a callback for non-critical connection errors
	// Multiple callbacks allowed
	OnError(func(error))

	// OnClose sets a callback for connection close. Error can be nil
	// Multiple callbacks allowed
	OnClose(func(error))

	// CloseChan returns channel for listener closing
	//
	// This method may not be used by end users,
	// it used by listener's connections to close properly
	CloseChan() <-chan error
}
