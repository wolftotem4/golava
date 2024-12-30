package env

import (
	"encoding/base64"
	"strings"
)

func Bool(value string) bool {
	return value == "true" || value == "1"
}

func Bytes(value string) ([]byte, error) {
	if strings.HasPrefix(value, "base64:") {
		return base64.StdEncoding.DecodeString(value[7:])
	}
	return []byte(value), nil
}
