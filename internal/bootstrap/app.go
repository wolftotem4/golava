package bootstrap

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava-core/golava"
	"github.com/wolftotem4/golava-core/hashing"
	"github.com/wolftotem4/golava-core/router"
	"github.com/wolftotem4/golava-core/session"
	"github.com/wolftotem4/golava/internal/app"
	_ "modernc.org/sqlite"
)

func InitApp(ctx context.Context) (*app.App, error) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var debug bool
	if debugVal := os.Getenv("APP_DEBUG"); debugVal == "true" || debugVal == "1" {
		debug = true
	}

	appKey, err := appKey()
	if err != nil {
		return nil, err
	}

	router, err := router.NewRouter("/")
	if err != nil {
		return nil, err
	}

	encrypter, err := initEncryption()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("sqlite", "db.sqlite")
	if err != nil {
		return nil, err
	}

	hasher := hashing.NewHasherManager()
	cookie := cookie.NewEncryptableCookieManager(initCookie(), encrypter)
	session := initSession(db)
	// session := initSession(&session.CookieSessionHandler{Cookie: cookie, Expiration: getSessionLifetime()})

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
		},
	}, nil
}

func initCookie() *cookie.CookieManager {
	path := os.Getenv("SESSION_PATH")
	if path == "" {
		path = "/"
	}

	var secure bool
	secureStr := os.Getenv("SESSION_SECURE_COOKIE")
	if secureStr == "true" || secureStr == "1" {
		secure = true
	}

	var sameSite http.SameSite
	switch os.Getenv("SESSION_SAME_SITE") {
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}

	return &cookie.CookieManager{
		Path:     path,
		Domain:   os.Getenv("SESSION_DOMAIN"),
		Secure:   secure,
		SameSite: sameSite,
	}
}

func initEncryption() (encryption.IEncrypter, error) {
	key, err := appKey()
	if err != nil {
		return nil, err
	}

	return encryption.NewEncrypter(key), nil
}

func appKey() ([]byte, error) {
	appKey := os.Getenv("APP_KEY")
	if appKey == "" {
		return nil, errors.New("APP_KEY is required")
	}

	if strings.HasPrefix(appKey, "base64:") {
		return base64.StdEncoding.DecodeString(appKey[7:])
	}

	return []byte(appKey), nil
}

func initSession(db *sqlx.DB) *session.SessionFactory {
	sessionName := os.Getenv("SESSION_NAME")
	if sessionName == "" {
		sessionName = "app_session"
	}

	httpOnly := true
	httpOnlyStr := os.Getenv("SESSION_HTTP_ONLY")
	if httpOnlyStr != "" && httpOnlyStr != "true" && httpOnlyStr != "1" {
		httpOnly = false
	}

	return &session.SessionFactory{
		Name:     sessionName,
		Lifetime: getSessionLifetime(),
		HttpOnly: httpOnly,
		Handler: &session.DatabaseSessionHandler{
			DB:         db.DB,
			DriverName: db.DriverName(),
		},
	}
}

func getSessionLifetime() time.Duration {
	lifeTime, _ := strconv.ParseInt(os.Getenv("SESSION_LIFETIME"), 10, 64)
	if lifeTime == 0 {
		lifeTime = 120
	}

	return time.Duration(lifeTime) * time.Minute
}
