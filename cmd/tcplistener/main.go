package main

import (
	"fmt"
	"github/HTTP-Protocol-From-TCP/internal/request"
	"log"
	"net"
)

func main() {

	f, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for {

		con, err := f.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection accepted")

		req, err := request.RequestFromReader(con)
		if err != nil {
			fmt.Println("error:", err)
			con.Close()
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		con.Close()

		fmt.Println("connection closed")
	}
}
