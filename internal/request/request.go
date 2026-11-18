package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState int

const (
	stateInitialized parserState = iota
	stateDone
)

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, 8)
	readToIndex := 0
	req := &Request{state: stateInitialized}

	for req.state != stateDone {
		if readToIndex >= len(buffer) {
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}

		n, readErr := reader.Read(buffer[readToIndex:])
		readToIndex += n

		consumed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[consumed:readToIndex])
		readToIndex -= consumed

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				if req.state != stateDone {
					return nil, fmt.Errorf("incomplete request...")
				}
				break
			}
			return nil, readErr
		}

	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	// requestLine = method request-target Http-version

	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return nil, 0, nil
	}

	line := string(data[:idx])

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line: %s", line)
	}
	method, target, version := parts[0], parts[1], parts[2]

	if method == "" {
		return nil, 0, fmt.Errorf("empty method")
	}
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, 0, fmt.Errorf("invalid method: %s", method)
		}
	}

	versionParts := strings.Split(version, "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return nil, 0, fmt.Errorf("invalid http version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: target,
		Method:        method,
	}, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = stateDone

		return n, nil

	case stateDone:
		return 0, fmt.Errorf("trying to read data in done state")

	default:
		return 0, fmt.Errorf("unknown state")
	}

}
