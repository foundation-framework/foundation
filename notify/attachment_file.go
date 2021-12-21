package notify

import (
	"io"
	"os"

	"github.com/foundation-framework/foundation/errors"
)

type fileAttachment struct {
	file *os.File
}

func NewFileAttachment(fileName string) Attachment {
	file, err := os.Open(fileName)
	if err != nil {
		errors.Panicf("failed to open file: %v", err)
	}

	return &fileAttachment{file: file}
}

func (f *fileAttachment) Name() string {
	return f.file.Name()
}

func (f *fileAttachment) Reader() io.Reader {
	return newOnlyReader(f.file)
}

func (f *fileAttachment) Reset() error {
	_, err := f.file.Seek(0, io.SeekStart)
	return err
}

func (f *fileAttachment) Close() error {
	return f.file.Close()
}
