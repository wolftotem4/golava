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
		i = instance.MustGetInstance(c)
		a = i.App.(*app.App)
	)

	guard := &generic.SessionGuard{
		Name:             "app",
		Session:          i.Session,
		Cookie:           i.Cookie,
		Hasher:           a.Hashing,
		RememberDuration: 400 * DAY,
		Provider: &db.SqlxUserProvider{
			Hasher:        a.Hashing,
			Table:         "users",
			DB:            a.DB,
			ConstructUser: func() auth.Authenticatable { return &generic.User{} },
		},
		RecallerIdMorph: auth.IntId,

		Request: c.Request,
	}

	i.Auth = guard

	err := guard.RestoreAuth(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}
