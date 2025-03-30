package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/livingpool/httpfromtcp/internal/request"
	"github.com/livingpool/httpfromtcp/internal/server"
)

const port = 42069

// Notice the sigChan code.
// This is a common pattern in Go for gracefully shutting down a server.
// Because server.Server returns immediately (it handles requests in the background in goroutines)
// if we exit main immediately, the server will just stop. We want to wait for a signal (like CTRL+C) before we stop the server.
func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Woopsie, my bad\n",
		}
	}
	w.Write([]byte("All good, frfr\n"))
	return nil
}
