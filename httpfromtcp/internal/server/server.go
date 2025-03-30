package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Port     int
	Listener net.Listener
	IsAlive  *atomic.Bool
}

func Serve(port int) (*Server, error) {
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
	response := "HTTP/1.1 200 OK\r\n" + // Status line
		"Content-Type: text/plain\r\n" + // Example header
		"\r\n" + // Blank line to separate headers from the body
		"Hello World!\n" // Body

	conn.Write([]byte(response))

	if err := conn.Close(); err != nil {
		log.Fatalf("error closing connection: %v", err)
	}
}
