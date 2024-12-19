package home

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava/binding"
)

func Login(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	data := t.Default(c)
	data["alert"], _ = instance.Session.Store.Get("alert")
	c.HTML(http.StatusOK, "home/login.tmpl", data)
}

func SubmitLogin(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	var form binding.Login
	if err := c.ShouldBind(&form); err != nil {
		instance.Session.Store.FlashInput(form)
		c.Error(err)
		return
	}

	statefulGuard, ok := instance.Auth.(auth.StatefulGuard)
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
		instance.Session.Store.Flash("alert", "Invalid credentials")
		instance.Redirector.Back(http.StatusSeeOther, "login")
		return
	}

	instance.Redirector.Intended(http.StatusSeeOther, "/")
}
