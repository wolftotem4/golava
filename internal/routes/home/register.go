package home

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/binding"
)

func RegisterView(c *gin.Context) {
	c.HTML(http.StatusOK, "home/register.tmpl", t.Default(c, t.PassFlash("alert-success", "alert-error")))
}

func Register(c *gin.Context) {
	var (
		i   = instance.MustGetInstance(c)
		app = i.App.(*app.App)
	)

	statefulGuard, ok := i.Auth.(auth.StatefulGuard)
	if !ok {
		c.Error(errors.New("auth does not implement StatefulGuard"))
		return
	}

	var form binding.Register
	if err := c.ShouldBind(&form); err != nil {
		i.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}

	row, err := app.DB.QueryContext(c, app.DB.Rebind("SELECT * FROM users WHERE username = ?"), form.Username)
	if err != nil {
		i.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}
	defer row.Close()

	if row.Next() {
		i.Session.Store.Flash("alert-error", "Username already exists")
		i.Session.Store.FlashInput(form)
		i.Redirector.Back(http.StatusSeeOther, "register")
		return
	}

	hash, err := app.Hashing.Make(form.Password)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := app.DB.ExecContext(c, app.DB.Rebind("INSERT INTO users (username, password) VALUES (?, ?)"), form.Username, hash)
	if err != nil {
		c.Error(err)
		return
	}

	userId, err := result.LastInsertId()
	if err != nil {
		c.Error(err)
		return
	}

	err = statefulGuard.LoginUsingID(c, userId, false)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusSeeOther, app.Router.URL("/").String())
}
