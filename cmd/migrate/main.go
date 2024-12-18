package main

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New("file://./database/migrations", "sqlite://./db.sqlite")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := m.Up(); err != nil {
		slog.Error(err.Error())
		return
	}
}
