package sockets

import (
	"context"
	"io"
	"net"
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
	HandleMsg(handlers ...Handler)

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

// Encoder describes connection encoding mechanism
type Encoder interface {
	// ResetReader resets reader used to decode data
	ResetReader(reader io.Reader)

	// ResetWriter resets writer used to encode data
	ResetWriter(writer io.Writer)

	// Flush enforces encoder to flush all buffered data
	Flush() error

	// ReadTopic reads message topic from underlying reader
	// Any encoding errors must be returned
	//
	// See Handler for more info about messaging system
	ReadTopic() (string, error)

	// ReadData reads message data from underlying reader
	// Any encoding returned by this method
	//
	// See Handler for more info about messaging system
	ReadData(data interface{}) error

	// WriteTopic writes message topic to underlying writer
	// Method will panic on any encoding errors
	//
	// See Handler for more info about messaging system
	WriteTopic(topic string) error

	// WriteData writes message data to underlying writer
	// Method will panic on any encoding errors
	//
	// See Handler for more info about messaging system
	WriteData(data interface{}) error
}

// Handler describes handler for incoming message on the connection
type Handler interface {
	// Context returns context that will be passed to Serve method
	Context() context.Context

	// Topic returns listening topic - base routing point
	//
	// Any topic format can be used, make sure it is consistent
	// with the other side of connection
	Topic() string

	// Model returns model used to decode messages with provided Topic
	// IMPORTANT!!! This method always must return pointer
	Model() interface{}

	// Serve used to serve message
	// Any data returned by this method will be sent back
	Serve(ctx context.Context, data interface{}) (string, interface{})
}

// HandlerFunc represents Handler calling function
type HandlerFunc func(ctx context.Context, data interface{}) (string, interface{})
