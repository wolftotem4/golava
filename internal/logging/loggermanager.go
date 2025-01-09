package logging

import (
	"log/slog"

	"github.com/pkg/errors"
)

type LoggerManager struct {
	logger map[string]*slog.Logger
}

func NewLoggerManager() *LoggerManager {
	return &LoggerManager{
		logger: make(map[string]*slog.Logger),
	}
}

func (lm *LoggerManager) Get(name string) (*slog.Logger, error) {
	logger, ok := lm.logger[name]
	if !ok {
		return nil, errors.New("logger not found")
	}
	return logger, nil
}

func (lm *LoggerManager) Set(name string, logger *slog.Logger) {
	lm.logger[name] = logger
}

func (lm *LoggerManager) MustGet(name string) *slog.Logger {
	logger, err := lm.Get(name)
	if err != nil {
		panic(err)
	}
	return logger
}
