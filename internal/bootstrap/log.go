package bootstrap

import (
	"log/slog"
	"os"

	"github.com/wolftotem4/golava/internal/env"
)

func InitLogger() error {
	if env.Bool(os.Getenv("APP_DEBUG")) {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	// logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	return err
	// }

	// slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(os.Stderr, logFile), nil)))

	return nil
}
