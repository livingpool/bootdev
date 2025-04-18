package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}

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
	case "/video":
		videoHandler(w, req)
	default:
		w.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(len(successHTML))
		h.Override("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(successHTML))
	}
}

// I recommend using netcat to test your chunked responses.
// Curl will abstract away the chunking for you, so you won't see your hex and cr and lf characters in your terminal if you use curl.
// I used this command to see my raw chunked response:
// echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069
// use curl --raw -v to view the whole thing
func proxyHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget

	url := "https://httpbin.org/" + strings.TrimPrefix(target, "/httpbin")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error connecting to httpbin.org: %v", err)
	}

	// the headers are only written at the start
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	h.Delete("Content-Length")
	w.WriteHeaders(h)

	buf := make([]byte, 1024)
	fullRespBody := make([]byte, 1024)
	fullRespBodyLen := 0

	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if _, err := w.WriteChunkedBodyDone(); err != nil {
					log.Fatalf("error writing chunked body done: %v", err)
				}
				hash := sha256.Sum256(fullRespBody[:fullRespBodyLen])
				headers := response.GetEmptyHeaders()
				headers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
				headers.Set("X-Content-Length", strconv.Itoa(fullRespBodyLen))
				w.WriteTrailers(headers)
				return
			}
			log.Fatalf("error reading from httpbin.org: %v", err)
		}

		if fullRespBodyLen+n >= cap(fullRespBody) {
			tempBuf := make([]byte, len(fullRespBody)*2)
			copy(tempBuf, fullRespBody)
			fullRespBody = tempBuf
		}
		copy(fullRespBody[fullRespBodyLen:], buf[:n])
		fullRespBodyLen += n

		fmt.Printf("read %d bytes from httpbin.org...\n", n)

		if _, err := w.WriteChunkedBody(buf[:n]); err != nil {
			log.Fatalf("error writing chunked body: %v", err)
		}
	}
}

// navigate to http://localhost:42069/video in your browser... does it work?
func videoHandler(w *response.Writer, req *request.Request) {
	f, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		log.Fatalf("error reading video file: %v", err)
	}

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(f))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(f)
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
