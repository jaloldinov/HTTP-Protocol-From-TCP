package request

import (
	"bytes"
	"errors"
	"fmt"
	"github/HTTP-Protocol-From-TCP/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	requestStateParsingHeaders
	requestStateParsingBody
)

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, 8)
	readToIndex := 0
	req := &Request{state: stateInitialized, Headers: headers.NewHeaders()}

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

func (r *Request) parseSingle(data []byte) (int, error) {
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
		r.state = requestStateParsingHeaders

		return n, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateParsingBody
		}

		return n, nil

	case requestStateParsingBody:
		contentLengthStr := r.Headers.Get("Content-Length")
		if contentLengthStr == "" {
			r.state = stateDone
			return 0, nil
		}

		contentLen, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length: %s", contentLengthStr)
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > contentLen {
			return 0, fmt.Errorf("body length are longer than content-length")
		}
		if len(r.Body) == contentLen {
			r.state = stateDone
		}
		return len(data), nil

	case stateDone:
		return 0, fmt.Errorf("trying to read data in done state")

	default:
		return 0, fmt.Errorf("unknown state")
	}

}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != stateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}
