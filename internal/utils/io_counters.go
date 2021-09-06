package utils

import "io"

type WriteCounter struct {
	writer io.Writer
	count  uint64
}

func (w *WriteCounter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *WriteCounter) Reset(writer io.Writer) {
	w.writer = writer
}

func (w *WriteCounter) Count() uint64 {
	return w.count
}

type ReadCounter struct {
	reader io.Reader
	count  uint64
}

func (w *ReadCounter) Read(data []byte) (int, error) {
	n, err := w.reader.Read(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *ReadCounter) ResetReader(reader io.Reader) {
	w.reader = reader
}

func (w *ReadCounter) Count() uint64 {
	return w.count
}
