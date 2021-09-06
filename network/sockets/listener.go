package sockets

// TODO: Documentation

type Listener interface {
	Listen(path string) error
	Close() error

	OnConnection(func(connection Connection))
	OnError(func(err error))
	OnClose(func(err error))
}
