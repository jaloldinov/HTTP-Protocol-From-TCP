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

	h[key] = strings.TrimSpace(value)

	return idx + 2, false, nil
}
