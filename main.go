package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		lines := getLinesChannel(con)
		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		buf := make([]byte, 8)
		var currentLine string
		for {
			n, err := f.Read(buf)

			if n > 0 {
				parts := strings.Split(string(buf[:n]), "\n")

				for i := 0; i < len(parts)-1; i++ {
					lines <- currentLine + parts[i]
					currentLine = ""
				}

				currentLine += parts[len(parts)-1]
			}

			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		if currentLine != "" {
			lines <- currentLine
		}
	}()

	return lines
}
