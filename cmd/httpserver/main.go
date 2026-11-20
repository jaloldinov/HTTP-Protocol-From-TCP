package main

import (
	"errors"
	"fmt"
	"github/HTTP-Protocol-From-TCP/internal/request"
	"github/HTTP-Protocol-From-TCP/internal/response"
	"github/HTTP-Protocol-From-TCP/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}
	status := response.StatusOK
	body := html200

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		body = html400

	case "/myproblem":
		status = response.StatusInternalServerError
		body = html500
	}

	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")

	w.WriteStatusLine(status)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))

}

const html400 = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const html500 = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const html200 = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

func proxyHandler(w *response.Writer, req *request.Request) {

	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	url := "https://httpbingo.org" + path

	resp, err := http.Get(url)
	if err != nil {
		h := response.GetDefaultHeaders(len(html500))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(h)
		w.WriteBody([]byte(html500))
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Delete("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(h)

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			if _, err := w.WriteChunkedBody(buf[:n]); err != nil {
				fmt.Println("error writing chunk:", err)
				break
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error reading response body:", err)
			break
		}
	}

	w.WriteChunkedBodyDone()
}
