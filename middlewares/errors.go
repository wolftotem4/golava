package middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/httputils"
	"github.com/wolftotem4/golava-core/httputils/csrf"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
)

func ErrorHandle(c *gin.Context) {
	instance := instance.MustGetInstance(c)

	c.Next()

	err := c.Errors.Last()
	if err == nil {
		return
	}

	if errors.Is(err, auth.ErrUnauthenticated) {
		unauthenticated(c, err, instance)
		return
	}

	if errors.Is(err, csrf.ErrTokenMismatch) {
		handleTokenMismatch(c, err)
		return
	}

	var validationError validator.ValidationErrors
	if errors.As(err, &validationError) {
		handleValidationError(c, validationError, instance)
		return
	}

	slog.ErrorContext(c, err.Error())

	if instance.App.Base().Debug {
		displayErrorMessage(c, http.StatusInternalServerError, err.Error())
	} else {
		displayErrorMessage(c, http.StatusInternalServerError, "Server Error")
	}
}

func handleValidationError(c *gin.Context, err validator.ValidationErrors, instance *instance.Instance) {
	var errors = make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Error()
	}

	if httputils.ExpectJson(c.Accepted) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Validation error", "errors": errors})
	} else {
		instance.Session.Store.Flash("errors", errors)
		instance.Redirector.Back(http.StatusSeeOther)
	}
}

func unauthenticated(c *gin.Context, err error, instance *instance.Instance) {
	if httputils.ExpectJson(c.Accepted) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	} else {
		instance.Redirector.Guest(http.StatusSeeOther, "login")
	}
}

func handleTokenMismatch(c *gin.Context, err error) {
	if httputils.ExpectJson(c.Accepted) {
		c.JSON(419, t.Default(c).Wrap(t.H{"message": err.Error()}))
	} else {
		c.HTML(419, "errors/419.tmpl", t.Default(c).Wrap(t.H{"message": err.Error()}))
	}
}

func displayErrorMessage(c *gin.Context, code int, message string) {
	if httputils.ExpectJson(c.Accepted) {
		c.JSON(code, t.Default(c).Wrap(t.H{"message": message}))
	} else {
		c.HTML(code, "errors/errors.tmpl", t.Default(c).Wrap(t.H{"message": message}))
	}
}
