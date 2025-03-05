package middlewares

import (
	"sync"

	"github.com/gin-gonic/gin"
	cookiemid "github.com/wolftotem4/golava-core/cookie/middleware"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/bootstrap"
)

type cookie struct {
	mu sync.Mutex
	gin.HandlerFunc
}

func (ci *cookie) Cookie(ctx *gin.Context) {
	if ci.HandlerFunc == nil {
		var (
			i = instance.MustGetInstance(ctx)
			a = i.App.(*app.App)
		)

		ci.makeCookieMiddleware(a)
	}

	ci.HandlerFunc(ctx)
}

func (c *cookie) makeCookieMiddleware(a *app.App) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.HandlerFunc == nil {
		factory := bootstrap.InitCookie(a.Encryption)
		c.HandlerFunc = cookiemid.CookieMiddleware(factory)
	}
}

var Cookie = (&cookie{}).Cookie
