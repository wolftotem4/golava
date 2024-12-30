package bootstrap

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/routing"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/env"
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

	session, err := initSession(db)
	if err != nil {
		return nil, err
	}

	translation, err := initTranslation(locale)
	if err != nil {
		return nil, err
	}

	hasher := hashing.NewHasherManager()
	cookie := cookie.NewEncryptableCookieManager(initCookie(), encrypter)

	return &app.App{
		DB: db,
		App: golava.App{
			Name:           os.Getenv("APP_NAME"),
			Debug:          debug,
			AppKey:         appKey,
			Router:         router,
			Cookie:         cookie,
			Encryption:     encrypter,
			Hashing:        hasher,
			SessionFactory: session,
			Translation:    translation,
			AppLocale:      locale,
		},
	}, nil
}
