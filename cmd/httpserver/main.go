package main

import (
	"github/HTTP-Protocol-From-TCP/internal/request"
	"github/HTTP-Protocol-From-TCP/internal/response"
	"github/HTTP-Protocol-From-TCP/internal/server"
	"io"
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

func handler(w io.Writer, req *request.Request) *server.HandleError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandleError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your proble is not my problem !\n",
		}

	case "/myproblem":
		return &server.HandleError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	}

	io.WriteString(w, "ALL good, frfr\n")
	return nil
}
