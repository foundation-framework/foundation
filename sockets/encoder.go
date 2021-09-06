package sockets

import (
	"io"
)

type Encoder interface {
	Reader(reader io.Reader)
	Writer(writer io.Writer)

	ReadTopic() (string, error)
	ReadData(data interface{}) error

	WriteTopic(topic string) error
	WriteData(data interface{}) error
}
