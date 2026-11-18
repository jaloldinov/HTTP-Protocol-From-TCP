package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil
}

func parseRequestLine(s string) (*RequestLine, error) {
	// requestLine = method request-target Http-version

	lines := strings.Split(s, "\r\n")
	line := lines[0]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", line)
	}
	method, target, version := parts[0], parts[1], parts[2]

	if method == "" {
		return nil, fmt.Errorf("empty method")
	}
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	versionParts := strings.Split(version, "/")
	if len(versionParts) != 2 || versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return nil, fmt.Errorf("invalid http version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: target,
		Method:        method,
	}, nil
}
