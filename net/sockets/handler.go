package sockets

import "context"

type HandlerBase interface {
	// Model returns model used to decode data
	// Must return pointer type
	Model() interface{}
}

//
// MessageHandler represents handler used to handle incoming messages
//
type MessageHandler interface {
	HandlerBase

	// Topic returns listening topic - base routing point
	//
	// Any topic format can be used, make sure it is consistent
	// with the other side of connection
	Topic() string

	// Serve used to serve message
	//
	// Any data returned by this method will be sent back
	// and trigger ReplyHandler (if present)
	Serve(data interface{}) interface{}
}

//
// MessageHandlerFunc represents MessageHandler.Serve function with context
// (This type used in handler implementation)
//
type MessageHandlerFunc func(ctx context.Context, data interface{}) interface{}

//
// MessageHandlerMiddleware represents simple middleware for handling message
//
type MessageHandlerMiddleware func(next MessageHandlerFunc) MessageHandlerFunc

//
// ReplyHandler represents handler used to handle replies from MessageHandler
//
type ReplyHandler interface {
	HandlerBase

	// Serve used to serve reply sent from MessageHandler
	Serve(data interface{})
}

//
// ReplyHandlerFunc represents ReplyHandler.Serve function with context
// (This type used in handler implementation)
//
type ReplyHandlerFunc func(ctx context.Context, data interface{})

//
// ReplyHandlerMiddleware represents simple middleware for handling reply
//
type ReplyHandlerMiddleware func(next ReplyHandlerFunc) ReplyHandlerFunc
