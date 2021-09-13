package transport

import "io"

// Encoder describes connection encoding mechanism
type Encoder interface {

	// ResetReader resets reader used to decode data
	ResetReader(reader io.Reader)

	// ResetWriter resets writer used to encode data
	ResetWriter(writer io.Writer)

	// Flush enforces encoder to flush all buffered data
	Flush() error

	// ReadTopic reads message topic from underlying reader
	// Any encoding errors must be returned
	//
	// See MessageHandler for more info about messaging system
	ReadTopic() (string, error)

	// ReadData reads message data from underlying reader
	// Any encoding returned by this method
	//
	// See MessageHandler for more info about messaging system
	ReadData(data interface{}) error

	// WriteTopic writes message topic to underlying writer
	// Method will panic on any encoding errors
	//
	// See MessageHandler for more info about messaging system
	WriteTopic(topic string) error

	// WriteData writes message data to underlying writer
	// Method will panic on any encoding errors
	//
	// See MessageHandler for more info about messaging system
	WriteData(data interface{}) error
}
