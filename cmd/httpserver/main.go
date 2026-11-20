package main

import (
	"github/HTTP-Protocol-From-TCP/internal/request"
	"github/HTTP-Protocol-From-TCP/internal/response"
	"github/HTTP-Protocol-From-TCP/internal/server"
	"log"
	"os"
	"os/signal"
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
