package sockets

import (
	"net/http"
)

// Server describes server for real-time communications
type Server interface {
	// Handler returns http.Handler for handling incoming connections
	Handler() http.Handler

	// HandleConn sets callback for an incoming connection
	// Multiple callbacks allowed
	HandleConn(func(Conn))

	// HandleError sets a callback for non-critical connection errors
	// Multiple callbacks allowed
	HandleError(func(error))

	// HandleClose sets a callback for listener closure
	// Error can be nil
	//
	// Multiple callbacks allowed
	HandleClose(func(error))
}
