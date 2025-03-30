package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/livingpool/httpfromtcp/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type Request struct {
	RequestLine  RequestLine
	requestState requestState
	Headers      headers.Headers
	Body         []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := &Request{
		requestState: requestStateInitialized,
		Headers:      headers.NewHeaders(),
		Body:         make([]byte, 0),
	}

	for request.requestState != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(io.EOF, err) {
				if request.requestState != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.requestState, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.requestState != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.requestState {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil { // something actually went wrong
			return 0, err
		}
		if n == 0 { // just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.requestState = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.requestState = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		// assume if no content-type header is present, there is no body
		if contentLength, exists := r.Headers.Get("Content-Length"); !exists {
			r.requestState = requestStateDone
			return 0, nil
		} else {
			r.Body = append(r.Body, data...)
			leng, err := strconv.Atoi(contentLength)
			if err != nil {
				return len(r.Body), fmt.Errorf("Content-Length is not integer, got=%s", contentLength)
			}
			if len(r.Body) > leng {
				return len(r.Body), fmt.Errorf("got more data then Content-Length, got=%d, expected:%d", len(r.Body), leng)
			}
			if len(r.Body) == leng {
				r.requestState = requestStateDone
			}
			return len(data), nil
		}
	case requestStateDone:
		return 0, fmt.Errorf("trying to read data from a requestStateDone requestState")
	default:
		return 0, fmt.Errorf("unknown requestState")
	}
}
