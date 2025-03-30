package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/livingpool/httpfromtcp/internal/request"
	"github.com/livingpool/httpfromtcp/internal/response"
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

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusBadRequest)
		h := response.GetDefaultHeaders(len(badReqHTML))
		h.Override("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(badReqHTML))
	case "/myproblem":
		w.WriteStatusLine(response.StatusInternalError)
		h := response.GetDefaultHeaders(len(internalErrorHTML))
		h.Override("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(internalErrorHTML))
	default:
		w.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(len(successHTML))
		h.Override("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(successHTML))
	}
}

const (
	successHTML = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

	badReqHTML = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

	internalErrorHTML = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
)
