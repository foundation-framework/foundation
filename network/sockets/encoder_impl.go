package sockets

import (
	"errors"
	"fmt"
	"io"

	"github.com/tinylib/msgp/msgp"
)

type msgpackEncoder struct {
	reader *msgp.Reader
	writer *msgp.Writer
}

func NewMsgpackEncoder() Encoder {
	return &msgpackEncoder{
		reader: msgp.NewReader(nil),
		writer: msgp.NewWriter(nil),
	}
}

func (e *msgpackEncoder) ResetReader(reader io.Reader) {
	e.reader.Reset(reader)
}

func (e *msgpackEncoder) ResetWriter(writer io.Writer) {
	e.writer.Reset(writer)
}

func (e *msgpackEncoder) Flush() error {
	return e.writer.Flush()
}

func (e *msgpackEncoder) ReadTopic() (string, error) {
	return e.reader.ReadString()
}

func (e *msgpackEncoder) ReadData(data interface{}) error {
	decodable, ok := data.(msgp.Decodable)
	if !ok {
		return errors.New("data is not msgp.Decodable")
	}

	if err := decodable.DecodeMsg(e.reader); err != nil {
		return err
	}

	return nil
}

func (e *msgpackEncoder) WriteTopic(topic string) error {
	return e.writer.WriteString(topic)
}

func (e *msgpackEncoder) WriteData(data interface{}) error {
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