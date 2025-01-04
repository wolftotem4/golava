package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/golava"
)

type App struct {
	golava.App

	Loggers *Loggers
	DB      *sqlx.DB
}

func (a *App) Base() *golava.App {
	return &a.App
}
