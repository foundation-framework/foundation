package notify

import "io"

type Attachment interface {
	Name() string
	Reader() io.Reader
	Reset() error
	Close() error
}
