package sockets

import (
	"context"
	"net"
)

// Conn describes a real-time connection
type Conn interface {
	// Accept method allows the connection to start reading messages
	Accept()

	// LocalAddr returns local endpoint address
	LocalAddr() net.Addr

	// RemoteAddr returns remote endpoint address
	RemoteAddr() net.Addr

	// BytesSent returns the total number of bytes sent
	BytesSent() uint64

	// BytesReceived returns the total number of bytes received
	BytesReceived() uint64

	// SetEncoder sets encoder for a connection (MessagePack Encoder used by default)
	//
	// You have to make sure that the same encoder is used
	// on the other side of the connection
	SetEncoder(encoder Encoder)

	// Write writers new message to the connection
	// Encoder used to encode message, use SetEncoder to change it
	Write(topic string, data interface{}, handler ...Handler) error

	// SetMessageHandlers sets handler for incoming message
	// Encoder used to decode message, use SetEncoder to change it
	//
	// Multiple handlers allowed
	SetMessageHandlers(handlers ...Handler)

	// RemoveMessageHandlers removes specifies message handlers
	RemoveMessageHandlers(handlers ...Handler)

	// OnError sets callback for non-critical connection errors
	//
	// Only one callback allowed, next calls will replace callback
	OnError(func(err error))

	// OnFatal sets callback for critical handler errors (usually panics)
	//
	// Only one callback allowed, next calls will replace callback
	OnFatal(func(topic string, data interface{}, msg interface{}))

	// OnClose sets callback for connection closure
	//
	// Multiple callback allowed
	OnClose(func(err error))

	// Close closes the connection
	Close(ctx context.Context) error
}
