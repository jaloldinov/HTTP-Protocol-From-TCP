package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLinesChannel(f)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
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
