package response

import (
	"fmt"
	"github/HTTP-Protocol-From-TCP/internal/headers"
	"io"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w, state: writerStateStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("status line must be written first!")
	}
	if err := WriteStatusLine(w.writer, statusCode); err != nil {
		return err
	}

	w.state = writerStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateHeaders {
		return fmt.Errorf("headers must be written after status line!")
	}

	if err := WriteHeaders(w.writer, h); err != nil {
		return err
	}
	w.state = writerStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("body must be written after headers!")
	}
	return w.writer.Write(p)
}
