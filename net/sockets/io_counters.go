package sockets

import "io"

type writeCounter struct {
	writer io.Writer
	count  uint64
}

func (w *writeCounter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *writeCounter) Reset(writer io.Writer) {
	w.writer = writer
}

func (w *writeCounter) Count() uint64 {
	return w.count
}

type readCounter struct {
	reader io.Reader
	count  uint64
}

func (w *readCounter) Read(data []byte) (int, error) {
	n, err := w.reader.Read(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *readCounter) ResetReader(reader io.Reader) {
	w.reader = reader
}

func (w *readCounter) Count() uint64 {
	return w.count
}
