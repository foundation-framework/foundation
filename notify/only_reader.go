package notify

import "io"

type onlyReader struct {
	inner io.Reader
}

func newOnlyReader(reader io.Reader) io.Reader {
	return &onlyReader{inner: reader}
}

func (o *onlyReader) Read(p []byte) (int, error) {
	return o.inner.Read(p)
}
