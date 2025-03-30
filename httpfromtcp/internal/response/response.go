package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/livingpool/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK            = StatusCode(200)
	StatusBadRequest    = StatusCode(400)
	StatusInternalError = StatusCode(500)
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func WriteBody(w io.Writer, body []byte) (int, error) {
	n, err := w.Write(body)
	return n, err
}
