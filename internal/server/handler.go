package server

import (
	"github/HTTP-Protocol-From-TCP/internal/request"
	"github/HTTP-Protocol-From-TCP/internal/response"
	"io"
)

type Handler func(w io.Writer, req *request.Request) *HandleError

type HandleError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he *HandleError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(he.Message)))
	w.Write([]byte(he.Message))
}
