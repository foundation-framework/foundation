package notify

import (
	"bytes"
	"io"
)

type memoryAttachment struct {
	name string
	data *bytes.Reader
}

func NewMemoryAttachment(name string, data []byte) Attachment {
	return &memoryAttachment{name: name, data: bytes.NewReader(data)}
}

func NewMemoryAttachmentString(name, data string) Attachment {
	return &memoryAttachment{name: name, data: bytes.NewReader([]byte(data))}
}

func (f *memoryAttachment) Name() string {
	return f.name
}

func (f *memoryAttachment) Reader() io.Reader {
	return f.data
}

func (f *memoryAttachment) Reset() error {
	_, err := f.data.Seek(0, io.SeekStart)

	// No errors expected
	return err
}

func (f *memoryAttachment) Close() error {
	// Do nothing
	return nil
}
