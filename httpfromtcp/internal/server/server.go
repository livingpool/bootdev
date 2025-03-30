package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/livingpool/httpfromtcp/internal/request"
	"github.com/livingpool/httpfromtcp/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	IsAlive  *atomic.Bool
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	state := &atomic.Bool{}
	state.Store(true)

	server := &Server{
		Port:     port,
		Listener: listener,
		IsAlive:  state,
		Handler:  handler,
	}

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.IsAlive.Store(false)
	err := s.Listener.Close()
	return err
}

func (s *Server) listen() {
	for {
		if s.IsAlive.Load() == false {
			return
		}
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}

		fmt.Println("connection accepted:", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("error reading request: %v", err)
	}

	buf := bytes.NewBuffer([]byte{})
	handlerError := s.Handler(buf, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	if err := response.WriteStatusLine(conn, 200); err != nil {
		log.Fatalf("error writing response status line: %v", err)
	}

	body := buf.Bytes()

	headers := response.GetDefaultHeaders(len(body))
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Fatalf("error writing respones headers: %v", err)
	}

	if _, err := response.WriteBody(conn, body); err != nil {
		log.Fatalf("error writing response body: %v", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatalf("error closing connection: %v", err)
	}
}
