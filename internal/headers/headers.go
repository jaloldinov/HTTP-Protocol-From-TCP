package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	line := strings.TrimSpace(string(data[:idx]))

	key, value, found := strings.Cut(line, ":")
	if !found {
		return 0, false, fmt.Errorf("invalid header: %s", line)
	}

	if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("invalid header key: %s", line)
	}

	if !isValidKey(key) {
		return 0, false, fmt.Errorf("invalid header key: %s", key)
	}

	key = strings.ToLower(key)
	value = strings.TrimSpace(value)

	if exist, ok := h[key]; ok {
		h[key] = exist + ", " + value
	} else {
		h[key] = value
	}

	return idx + 2, false, nil
}

func isValidKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	for _, c := range key {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
			continue
		}

		if strings.ContainsRune("!#$%&'*+-.^_`|~", c) {
			continue
		}

		return false
	}

	return true
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Delete(key string) {
	delete(h, strings.ToLower(key))
}
