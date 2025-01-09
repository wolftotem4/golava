package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava/internal/logging"
)

type App struct {
	golava.App

	L  *logging.LoggerManager
	DB *sqlx.DB
}

func (a *App) Base() *golava.App {
	return &a.App
}
