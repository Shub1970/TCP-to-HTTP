package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

// function should only return done=true when the data starts with a CRLF, which can't happen when it finds a new key/value pair.

const CRLF = "\r\n"

func fieldNameVerifier(name string) bool {
	for _, ch := range name {
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' {
			continue
		}

		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch == '!' || ch == '#' || ch == '$' || ch == '%' || ch == '&' || ch == '\'' || ch == '*' || ch == '+' || ch == '-' || ch == '.' || ch == '^' || ch == '_' || ch == '`' || ch == '|' || ch == '~' {
			continue
		}
		return false
	}
	return true
}

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(key string) (string, bool) {
	lowerKey := strings.ToLower(key)
	if val, ok := h[lowerKey]; ok {
		return val, ok
	}
	return "", false
}

func (h Headers) Set(key string, newValue string) {
	if value, ok := h[key]; ok {
		h[key] = fmt.Sprintf("%s,%s", value, newValue)
	} else {
		h[key] = newValue
	}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	readLen := 0
	fieldLinePairIndex := bytes.Index(data, []byte(CRLF))
	if fieldLinePairIndex == -1 {
		return 0, false, nil
	}
	if fieldLinePairIndex == 0 {
		return 0, true, nil
	}
	readLen += fieldLinePairIndex + len(CRLF)
	// split around colon
	properFieldLine := bytes.TrimSpace(data[:fieldLinePairIndex])
	headerParts := bytes.SplitN(properFieldLine, []byte(":"), 2)
	if len(headerParts) != 2 {
		return 0, false, fmt.Errorf("invalid header line: %s", string(properFieldLine))
	}
	key := headerParts[0]
	// there should be no space b/w the colon and field name
	if bytes.HasSuffix(key, []byte(" ")) {
		return 0, false, fmt.Errorf("invalid field name: %s", string(key))
	}
	key = bytes.ToLower(bytes.TrimSpace(key))
	value := bytes.ToLower(bytes.TrimSpace(headerParts[1]))
	if !fieldNameVerifier(string(key)) {
		return 0, false, fmt.Errorf("invalid field name verified: %s", string(key))
	}
	h.Set(string(key), string(value))

	return readLen, false, nil
}
