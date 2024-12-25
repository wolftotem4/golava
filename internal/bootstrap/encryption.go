package bootstrap

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/wolftotem4/golava-core/encryption"
)

func initEncryption() (encryption.IEncrypter, error) {
	key, err := appKey()
	if err != nil {
		return nil, err
	}

	return encryption.NewEncrypter(key), nil
}

func appKey() ([]byte, error) {
	appKey := os.Getenv("APP_KEY")
	if appKey == "" {
		return nil, errors.New("APP_KEY is required")
	}

	if strings.HasPrefix(appKey, "base64:") {
		return base64.StdEncoding.DecodeString(appKey[7:])
	}

	return []byte(appKey), nil
}
