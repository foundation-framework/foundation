package sockets

import "io"

// Encoder describes connection encoding mechanism
type Encoder interface {
	// ResetReader resets reader used to decode data
	ResetReader(reader io.Reader)

	// ResetWriter resets writer used to encode data
	ResetWriter(writer io.Writer)

	// Flush enforces encoder to flush all buffered data
	Flush() error

	// ReadString reads a string from an underlying reader
	// Any encoding errors must be returned
	ReadString() (string, error)

	// ReadData reads message data from underlying reader
	// Any encoding returned by this method
	ReadData(data interface{}) error

	// WriteString writes a string to an underlying writer
	// Method will panic on any encoding errors
	WriteString(content string) error

	// WriteData writes message data to underlying writer
	// Method will panic on any encoding errors
	WriteData(data interface{}) error
}
