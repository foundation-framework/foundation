package sockets

import (
	"errors"
	"fmt"
	"io"

	"github.com/tinylib/msgp/msgp"
)

type msgpEncoder struct {
	reader *msgp.Reader
	writer *msgp.Writer
}

func NewMsgpEncoder() Encoder {
	return &msgpEncoder{
		reader: msgp.NewReader(nil),
		writer: msgp.NewWriter(nil),
	}
}

func (e *msgpEncoder) Reader(reader io.Reader) {
	e.reader.Reset(reader)
}

func (e *msgpEncoder) Writer(writer io.Writer) {
	e.writer.Reset(writer)
}

func (e *msgpEncoder) ReadTopic() (string, error) {
	return e.reader.ReadString()
}

func (e *msgpEncoder) ReadData(data interface{}) error {
	decodable, ok := data.(msgp.Decodable)
	if !ok {
		return errors.New("data is not msgp.Decodable")
	}

	err := decodable.DecodeMsg(e.reader)
	if _, ok := err.(msgp.Error); ok {
		panic(fmt.Errorf("unexpected deconding error: %s", err.Error()))
	}

	return err
}

func (e *msgpEncoder) WriteTopic(topic string) error {
	return e.writer.WriteString(topic)
}

func (e *msgpEncoder) WriteData(data interface{}) error {
	encodable, ok := data.(msgp.Encodable)
	if !ok {
		return errors.New("data is not msgp.Encodable")
	}

	err := encodable.EncodeMsg(e.writer)
	if _, ok := err.(msgp.Error); ok {
		panic(fmt.Errorf("unexpected encoding error: %s", err.Error()))
	}

	return err
}
