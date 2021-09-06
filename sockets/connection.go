package sockets

import (
	"net"
)

// TODO: Documentation

type Connection interface {
	SetEncoder(encoder Encoder)

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Write(topic string, data interface{}) error
	Close() error

	BytesSent() uint64
	BytesReceived() uint64

	OnMessage(handler Handler)
	OnError(func(err error))
	OnClose(func(err error))
}
