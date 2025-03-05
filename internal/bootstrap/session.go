package bootstrap

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wolftotem4/golava-core/session"
	sess "github.com/wolftotem4/golava-core/session/sqlite"
	"github.com/wolftotem4/golava/internal/env"
)

func InitSession(db *sqlx.DB, table string) (*session.SessionFactory, error) {
	handler := sess.NewSqliteSessionHandler(db.DB, table)

	return &session.SessionFactory{
		Name:     env.Get("SESSION_COOKIE", "app_session"),
		Lifetime: getSessionLifetime(),
		HttpOnly: env.Bool(env.Get("SESSION_HTTP_ONLY", "true")),
		Handler:  handler,
	}, nil
}

func getSessionLifetime() time.Duration {
	lifeTime := env.Int(os.Getenv("SESSION_LIFETIME"))
	if lifeTime == 0 {
		lifeTime = 120
	}

	return time.Duration(lifeTime) * time.Minute
}
