package bootstrap

import (
	"log/slog"
	"os"

	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/env"
	"github.com/wolftotem4/golava/internal/logging"
)

func InitDefaultLogger() error {
	level := getDefaultLogLevel()
	slog.SetLogLoggerLevel(level)
	opts := slog.HandlerOptions{Level: level}

	handler, err := logging.GetHandler(env.Get("APP_LOG_SINK", "console"), &opts)
	if err != nil {
		return err
	}

	slog.SetDefault(slog.New(handler))

	return nil
}

func initLoggers() (*app.Loggers, error) {
	opts := &slog.HandlerOptions{Level: getDefaultLogLevel()}

	request, err := logging.GetHandler(env.Get("REQUEST_LOG_SINK", "console"), opts)
	if err != nil {
		return nil, err
	}

	return &app.Loggers{
		Request: slog.New(request),
	}, nil
}

func getDefaultLogLevel() slog.Level {
	if env.Bool(os.Getenv("APP_DEBUG")) {
		return slog.LevelDebug
	} else {
		return slog.LevelInfo
	}
}
