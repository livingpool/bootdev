package server

import (
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

type Handler func(w *response.Writer, req *request.Request)

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

	writer := response.NewResponseWriter(conn)
	s.Handler(writer, req)

	if err := conn.Close(); err != nil {
		log.Fatalf("error closing connection: %v", err)
	}
}
