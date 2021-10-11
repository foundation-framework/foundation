package websocket

import "io"

type writerCounter struct {
	writer io.Writer
	count  uint64
}

func (w *writerCounter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *writerCounter) Reset(writer io.Writer) {
	w.writer = writer
}

func (w *writerCounter) Count() uint64 {
	return w.count
}

type readerCounter struct {
	reader io.Reader
	count  uint64
}

func (w *readerCounter) Read(data []byte) (int, error) {
	n, err := w.reader.Read(data)
	if err != nil {
		return 0, err
	}

	w.count += uint64(n)
	return n, nil
}

func (w *readerCounter) ResetReader(reader io.Reader) {
	w.reader = reader
}

func (w *readerCounter) Count() uint64 {
	return w.count
}
