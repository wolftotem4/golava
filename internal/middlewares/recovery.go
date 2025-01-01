package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava/internal/helper"
)

func Recovery(appDebug bool) gin.HandlerFunc {
	return gin.RecoveryWithWriter(emptyWriter{}, func(c *gin.Context, err interface{}) {
		stacktrace := string(debug.Stack())

		slog.ErrorContext(c, fmt.Sprintf("%v", err), "stacktrace", stacktrace)

		var message string
		if appDebug {
			message = fmt.Sprintf("%v", err)
		} else {
			i := instance.MustGetInstance(c)
			message, _ = helper.GetTranslator(i, true).T("error.server_error")
		}

		if utils.ExpectJson(c.GetHeader("Accept")) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": message})
		} else if appDebug {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%s\n\n%s", message, stacktrace))
		} else {
			c.HTML(http.StatusInternalServerError, "errors/errors.tmpl", t.Default(c).Wrap(t.H{"message": message}))
		}
	})
}

type emptyWriter struct{}

func (emptyWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
