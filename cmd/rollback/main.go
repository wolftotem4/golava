package main

import (
	"flag"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var step int

func init() {
	flag.IntVar(&step, "step", 0, "step")
}

func main() {
	flag.Parse()

	if step < 0 {
		slog.Error("step must not be negative")
		return
	}

	m, err := migrate.New("file://./database/migrations", "sqlite://./db.sqlite")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if step > 0 {
		if err := m.Steps(-step); err != nil {
			slog.Error(err.Error())
			return
		}
		return
	}

	if err := m.Down(); err != nil {
		slog.Error(err.Error())
		return
	}
}
