package home

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
)

type LoginForm struct {
	Username string `json:"username" form:"username" mod:"trim" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Remember bool   `json:"remember" form:"remember"`
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "home/login.tmpl", t.Default(c, t.PassFlash("alert-success", "alert-error")))
}

func SubmitLogin(c *gin.Context) {
	i := instance.MustGetInstance(c)

	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		i.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}

	statefulGuard, ok := i.Auth.(auth.StatefulGuard)
	if !ok {
		c.Error(errors.New("auth does not implement StatefulGuard"))
		return
	}

	credentials := map[string]interface{}{
		"username": form.Username,
		"password": form.Password,
	}

	valid, err := statefulGuard.Attempt(c, credentials, form.Remember)
	if err != nil {
		c.Error(err)
		return
	} else if !valid {
		i.Session.Store.Flash("alert-error", "Invalid credentials")
		i.Redirector.Back(http.StatusSeeOther, "login")
		return
	}

	i.Redirector.Intended(http.StatusSeeOther, "/")
}
