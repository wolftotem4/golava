package bootstrap

import (
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/session"
	sess "github.com/wolftotem4/golava-core/session/sqlite"
)

func initSession(db *sqlx.DB) (*session.SessionFactory, error) {
	sessionName := os.Getenv("SESSION_COOKIE")
	if sessionName == "" {
		sessionName = "app_session"
	}

	httpOnly := true
	httpOnlyStr := os.Getenv("SESSION_HTTP_ONLY")
	if httpOnlyStr != "" && httpOnlyStr != "true" && httpOnlyStr != "1" {
		httpOnly = false
	}

	handler := &sess.SqliteSessionHandler{DB: db.DB}

	return &session.SessionFactory{
		Name:     sessionName,
		Lifetime: getSessionLifetime(),
		HttpOnly: httpOnly,
		Handler:  handler,
	}, nil
}

func getSessionLifetime() time.Duration {
	lifeTime, _ := strconv.ParseInt(os.Getenv("SESSION_LIFETIME"), 10, 64)
	if lifeTime == 0 {
		lifeTime = 120
	}

	return time.Duration(lifeTime) * time.Minute
}
