package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava/binding"
	"github.com/wolftotem4/golava/internal/app"
)

func RegisterView(c *gin.Context) {
	c.HTML(http.StatusOK, "home/register.tmpl", t.Default(c))
}

func Register(c *gin.Context) {
	var (
		instance = instance.MustGetInstance(c)
		app      = instance.App.(*app.App)
	)

	var form binding.Register
	if err := c.ShouldBind(&form); err != nil {
		instance.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}

	row, err := app.DB.QueryContext(c, "SELECT * FROM users WHERE name = ?", form.Username)
	if err != nil {
		instance.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}
	defer row.Close()

	if row.Next() {
		instance.Session.Store.Flash("alert", "Username already exists")
		instance.Session.Store.FlashInput(form)
		instance.Redirector.Back(http.StatusSeeOther, "register")
		return
	}

	hash, err := app.Hashing.Make(form.Password)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = app.DB.ExecContext(c, "INSERT INTO users (name, password) VALUES (?, ?)", form.Username, hash)
	if err != nil {
		c.Error(err)
		return
	}

	instance.Session.Store.Flash("alert", "You have been registered successfully")
	c.Redirect(http.StatusSeeOther, app.Router.URL("/").String())
}
