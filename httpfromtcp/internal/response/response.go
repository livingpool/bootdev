package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/livingpool/httpfromtcp/internal/headers"
)

type StatusCode int
type writerState int

const (
	StatusOK            = StatusCode(200)
	StatusBadRequest    = StatusCode(400)
	StatusInternalError = StatusCode(500)
)

const (
	writingStatusLine writerState = iota
	writingHeaders
	writingBody
	writingDone
)

type Writer struct {
	stream      io.Writer
	writerState writerState
}

func NewResponseWriter(stream io.Writer) *Writer {
	return &Writer{
		stream:      stream,
		writerState: writingStatusLine,
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.stream.Write(p)
	return n, err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writingStatusLine {
		return fmt.Errorf("state is not writingStatusLine")
	}

	reasonPhrase := ""

	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalError:
		reasonPhrase = "Internal Server Error"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	w.writerState = writingHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writingHeaders {
		return fmt.Errorf("state is not writingHeaders")
	}

	for k, v := range headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}

	w.writerState = writingBody
	_, err := w.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.writerState != writingBody {
		return 0, fmt.Errorf("state is not writingBody")
	}

	w.writerState = writingDone
	n, err := w.Write(body)
	return n, err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}
