package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava/internal/env"
)

func main() {
	err := godotenv.Load(env.DotEnvFile)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Generating new key...")
	key, err := encryption.GenerateKey()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = env.SetKeyInEnvironmentFile("APP_KEY", fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString(key)))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Key generated successfully")
}
