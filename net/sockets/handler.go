package sockets

import "context"

// Handler describes handler for incoming message on the connection
type Handler interface {
	// Context called first on new message and used to construct context
	Context(caller func(ctx context.Context))

	// Topic returns listening topic - base routing point
	//
	// Any topic format can be used, make sure it is consistent
	// with the other side of connection
	//
	// IMPORTANT!!! Reply handlers must have empty topic
	Topic() string

	// Model returns model used to decode messages with provided Topic
	// IMPORTANT!!! This method always must return pointer
	Model() interface{}

	// Serve used to serve message
	//
	// Any data returned by this method will be sent back
	// and trigger reply Handler passed to Write method
	Serve(ctx context.Context, data interface{}) interface{}
}

// HandlerFunc represents serving function for standard message handler
type HandlerFunc func(ctx context.Context, data interface{}) interface{}

// ReplyHandlerFunc represents serving function for reply handler
type ReplyHandlerFunc func(ctx context.Context, data interface{})
