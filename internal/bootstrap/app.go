package bootstrap

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/routing"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/env"
	"github.com/wolftotem4/golava/internal/logging"
	_ "modernc.org/sqlite"
)

func InitApp(ctx context.Context) (*app.App, error) {
	locale := "en"

	debug := env.Bool(os.Getenv("APP_DEBUG"))

	appKey, err := appKey()
	if err != nil {
		return nil, err
	}

	router, err := routing.NewRouter("/")
	if err != nil {
		return nil, err
	}
	router.BaseURL, _ = url.Parse(os.Getenv("BASE_URL"))

	encrypter, err := initEncryption()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	translation, err := initTranslation(locale)
	if err != nil {
		return nil, err
	}

	return &app.App{
		L:  logging.NewLoggerManager(),
		DB: db,
		App: golava.App{
			Name:        os.Getenv("APP_NAME"),
			Debug:       debug,
			AppKey:      appKey,
			Router:      router,
			Encryption:  encrypter,
			Hashing:     hashing.NewHasherManager(),
			Translation: translation,
			AppLocale:   locale,
		},
	}, nil
}
