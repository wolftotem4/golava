package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var step int
var reset bool

func init() {
	flag.IntVar(&step, "step", 0, "step")
	flag.BoolVar(&reset, "reset", false, "reset")
}

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var (
		driverName = os.Getenv("DB_DRIVER")
		dsn        = os.Getenv("DB_DSN")
	)

	driver, err := driver(driverName, dsn)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://./database/migrations", driverName, driver)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	m.Log = new(logger)

	if step != 0 {
		if err := m.Steps(step); err != nil {
			slog.Error(err.Error())
			return
		}
	} else if reset {
		if err := m.Down(); err != nil {
			slog.Error(err.Error())
			return
		}
		return
	} else {
		if err := m.Up(); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func driver(driverName string, dataSourceName string) (database.Driver, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	switch driverName {
	case "mysql":
		return mysql.WithInstance(db, &mysql.Config{})
	case "sqlite":
		return sqlite.WithInstance(db, &sqlite.Config{})
	case "postgres":
		return postgres.WithInstance(db, &postgres.Config{})
	default:
		return nil, fmt.Errorf("unknown driver: %s", driverName)
	}
}

type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func (l *logger) Verbose() bool {
	return false
}
