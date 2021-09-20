package sockets

import (
	"context"
	"net"
)

// Conn describes a real-time connection
type Conn interface {

	// SetEncoder sets encoder for a connection (MessagePack Encoder used by default)
	//
	// You have to make sure that the same encoder is used
	// on the other side of the connection
	SetEncoder(encoder Encoder)

	// LocalAddr returns local endpoint address
	LocalAddr() net.Addr

	// RemoteAddr returns remote endpoint address
	RemoteAddr() net.Addr

	// Write writers new message to the connection
	// Encoder used to encode message, use SetEncoder to change it
	Write(topic string, data interface{}) error

	// Close closes the connection
	Close(ctx context.Context) error

	// BytesSent returns the total number of bytes sent
	BytesSent() uint64

	// BytesReceived returns the total number of bytes received
	BytesReceived() uint64

	// HandleMsg sets handler for incoming message
	// Encoder used to decode message, use SetEncoder to change it
	//
	// Multiple handlers allowed
	HandleMsg(handler Handler)

	// HandleError sets callback for non-critical connection errors
	// Multiple callbacks allowed
	HandleError(func(err error))

	// HandleFatal sets callback for critical handler errors (usually panics)
	// Only one callback allowed, next calls will replace callback
	HandleFatal(func(topic string, data interface{}, msg interface{}))

	// HandleClose sets callback for connection closure
	// Multiple callbacks allowed
	HandleClose(func(err error))
}