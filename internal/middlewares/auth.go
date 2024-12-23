package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/auth/generic"
	"github.com/wolftotem4/golava-core/instance"
	db "github.com/wolftotem4/golava-db-sqlx"
	"github.com/wolftotem4/golava/internal/app"
)

func WebAuth(c *gin.Context) {
	const (
		DAY = 24 * time.Hour
	)

	var (
		instance = instance.MustGetInstance(c)
		app      = instance.App.(*app.App)
	)

	guard := &generic.SessionGuard{
		Name:             "app",
		Session:          instance.Session,
		Cookie:           app.Cookie,
		RememberDuration: 400 * DAY,
		Provider: &db.SqlxUserProvider{
			Hasher:        app.Hashing,
			Table:         "users",
			DB:            app.DB,
			ConstructUser: func() auth.Authenticatable { return &generic.User{} },
		},

		Request: c.Request,
	}

	instance.Auth = guard

	err := guard.RestoreAuth(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}
