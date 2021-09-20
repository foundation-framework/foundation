package sockets

import "context"

// Handler describes handler for incoming message on the connection
type Handler interface {

	// Topic returns listening topic - base routing point
	//
	// Any topic format can be used, make sure it is consistent
	// with the other side of connection
	Topic() string

	// Model returns model used to decode messages with provided Topic
	// IMPORTANT!!! This method always must return pointer
	Model() interface{}

	// Serve used to serve message
	// Any data returned by this method will be sent back using the same topic
	Serve(ctx context.Context, data interface{}) interface{}
}