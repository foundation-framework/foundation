package sockets

import "context"

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
	// and trigger ReplyHandler
	Serve(ctx context.Context, data interface{}) interface{}
}

// HandlerFunc represents Handler.Serve function
type HandlerFunc func(ctx context.Context, data interface{}) interface{}
