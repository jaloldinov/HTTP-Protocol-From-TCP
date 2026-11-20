package server

import (
	"fmt"
	"github/HTTP-Protocol-From-TCP/internal/request"
	"github/HTTP-Protocol-From-TCP/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: listener, handler: handler}
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		body := []byte(err.Error())
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(req.Body)
		return
	}

	s.handler(w, req)
}
