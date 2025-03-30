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

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 { // headers are done, consume the CRLF
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("space after field name: %s", key)
	}

	val := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}

	h.Set(key, string(val))

	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	existingVal, exists := h[key]
	if exists {
		value = existingVal + ", " + value
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	val, exists := h[key]
	return val, exists
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

// RFC 9110 5.1 and 5.6.2
var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !(c >= 'A' && c <= 'Z' ||
			c >= 'a' && c <= 'z' ||
			c >= '0' && c <= '9' ||
			c == '-') {
			return false
		}
	}
	return true
}
