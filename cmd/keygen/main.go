package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wolftotem4/golava-core/encryption"
)

const dotEnvFile = ".env"

func main() {
	err := godotenv.Load(dotEnvFile)
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

	err = setKeyInEnvironmentFile(fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString(key)))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("Key generated successfully")
}

func setKeyInEnvironmentFile(key string) error {
	old, err := os.ReadFile(dotEnvFile)
	if err != nil {
		return err
	}

	env := string(old)
	oldKey := os.Getenv("APP_KEY")

	if strings.Contains(env, fmt.Sprintf("APP_KEY=%s", oldKey)) {
		newEnv := strings.ReplaceAll(env, fmt.Sprintf("APP_KEY=%s", oldKey), fmt.Sprintf("APP_KEY=%s", key))

		return os.WriteFile(dotEnvFile, []byte(newEnv), 0644)
	} else {
		f, err := os.OpenFile(dotEnvFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(fmt.Sprintf("\nAPP_KEY=%s\n", key))
		return err
	}
}
