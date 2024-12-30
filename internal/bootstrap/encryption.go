package bootstrap

import (
	"errors"
	"os"

	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava/internal/env"
)

func initEncryption() (encryption.IEncrypter, error) {
	key, err := appKey()
	if err != nil {
		return nil, err
	}

	return encryption.NewEncrypter(key), nil
}

func appKey() ([]byte, error) {
	appKey, err := env.Bytes(os.Getenv("APP_KEY"))
	if err != nil {
		return nil, err
	}

	if len(appKey) == 0 {
		return nil, errors.New("APP_KEY is required")
	}

	return appKey, nil
}
