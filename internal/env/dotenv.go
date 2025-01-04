package env

import (
	"encoding/base64"
	"os"
	"strconv"
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

func Int(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}

func Get(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
