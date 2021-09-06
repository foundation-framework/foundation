package sockets

import (
	"io"
)

type Encoder interface {
	ResetReader(reader io.Reader)
	ResetWriter(writer io.Writer)
	Flush() error

	ReadTopic() (string, error)
	ReadData(data interface{}) error

	WriteTopic(topic string) error
	WriteData(data interface{}) error
}
