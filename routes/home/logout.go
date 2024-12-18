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
		instance = instance.MustGetInstance(ctx)
		app      = instance.App.(*app.App)
	)

	statefulGuard, ok := instance.Auth.(auth.StatefulGuard)
	if !ok {
		ctx.Error(errors.New("auth does not implement StatefulGuard"))
		return
	}

	err := statefulGuard.Logout(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	instance.Session.Store.Invalidate(ctx)
	instance.Session.Store.RegenerateToken()

	app.Cookie.Encryption().Set(
		instance.Session.GetMigrateName(),
		instance.Session.Store.ID,
		cookie.WithMaxAge(int(instance.Session.Lifetime)),
	)

	instance.Redirector.Redirect(http.StatusSeeOther, "/")
}
