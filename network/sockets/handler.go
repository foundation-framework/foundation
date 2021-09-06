package sockets

// TODO: Documentation

type Handler interface {
	Topic() string
	Model() interface{}

	Serve(data interface{}) interface{}
}
