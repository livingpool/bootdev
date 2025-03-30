package server

import (
	"io"

	"github.com/livingpool/httpfromtcp/internal/request"
	"github.com/livingpool/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	headers := response.GetDefaultHeaders(len(he.Message))
	response.WriteHeaders(w, headers)
	response.WriteBody(w, []byte(he.Message))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
