package sockets

import (
	"net/http"
)

// Server describes server for real-time communications
type Server interface {
	// Handler returns http.Handler for handling incoming connections
	Handler() http.Handler

	// OnConn sets callback for an incoming connection
	//
	// Only one callback allowed, next calls will replace callback
	OnConn(func(Conn, http.Header))

	// OnError sets a callback for non-critical connection errors
	//
	// Only one callback allowed, next calls will replace callback
	OnError(func(error))
}
