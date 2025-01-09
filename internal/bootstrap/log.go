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

func InitLogger(a *app.App, settings ...func(a *app.App) error) error {
	for _, setting := range settings {
		if err := setting(a); err != nil {
			return err
		}
	}

	return nil
}

func getDefaultLogLevel() slog.Level {
	if env.Bool(os.Getenv("APP_DEBUG")) {
		return slog.LevelDebug
	} else {
		return slog.LevelInfo
	}
}

func Logger(name string, logger *slog.Logger) func(a *app.App) error {
	return func(a *app.App) error {
		a.L.Set(name, logger)
		return nil
	}
}

func LoggerSink(name, sink string) func(a *app.App) error {
	return func(a *app.App) error {
		opts := &slog.HandlerOptions{Level: getDefaultLogLevel()}

		logger, err := logging.GetHandler(sink, opts)
		if err != nil {
			return err
		}

		a.L.Set(name, slog.New(logger))

		return nil
	}
}
