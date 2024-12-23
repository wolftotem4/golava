package home

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava/internal/app"
)

func SubmitLogout(ctx *gin.Context) {
	var (
		i   = instance.MustGetInstance(ctx)
		app = i.App.(*app.App)
	)

	statefulGuard, ok := i.Auth.(auth.StatefulGuard)
	if !ok {
		ctx.Error(errors.New("auth does not implement StatefulGuard"))
		return
	}

	err := statefulGuard.Logout(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	i.Session.Store.Invalidate(ctx)
	i.Session.Store.RegenerateToken()

	app.Cookie.Encryption().Set(
		i.Session.GetMigrateName(),
		i.Session.Store.ID,
		cookie.WithMaxAge(int(i.Session.Lifetime)),
	)

	i.Redirector.Redirect(http.StatusSeeOther, "/")
}
