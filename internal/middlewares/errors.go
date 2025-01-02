package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/wolftotem4/golava-core/auth"
	"github.com/wolftotem4/golava-core/http/csrf"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava/internal/helper"
	"github.com/wolftotem4/golava/internal/ratelimit"
)

type stackTraceableError interface {
	error
	StackTrace() errors.StackTrace
}

func ErrorHandle(c *gin.Context) {
	i := instance.MustGetInstance(c)

	c.Writer = &WriterMonitor{ResponseWriter: c.Writer}

	c.Next()

	err := c.Errors.Last()
	if err == nil {
		return
	}

	if errors.Is(err, auth.ErrUnauthenticated) {
		unauthenticated(c, err, i)
		return
	}

	if errors.Is(err, csrf.ErrTokenMismatch) {
		handleTokenMismatch(c, err)
		return
	}

	if errors.Is(err, ratelimit.ErrTooManyAttempts) {
		displayErrorMessage(c, http.StatusTooManyRequests, err.Error())
		return
	}

	var validationError validator.ValidationErrors
	if errors.As(err, &validationError) {
		handleValidationError(c, validationError, i)
		return
	}

	var (
		stackTracable bool
		stackTraceErr stackTraceableError
	)
	if stackTracable = errors.As(err, &stackTraceErr); stackTracable {
		slog.ErrorContext(c, fmt.Sprintf("%+v", stackTraceErr))
	} else {
		slog.ErrorContext(c, err.Error())
	}

	if i.App.Base().Debug {
		if stackTracable {
			displayErrorWithStackTrace(c, http.StatusInternalServerError, stackTraceErr)
		} else {
			displayErrorMessage(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		msg, _ := helper.GetTranslator(i, true).T("error.server_error")
		displayErrorMessage(c, http.StatusInternalServerError, msg)
	}
}

func handleValidationError(c *gin.Context, err validator.ValidationErrors, i *instance.Instance) {
	trans := i.GetUserPreferredTranslator() // it has to be original translator, we can't use `helper.GetTranslator` here

	var errors = make(map[string]string)
	for _, e := range err {
		errors[e.Field()] = e.Translate(trans)
	}

	if utils.ExpectJson(c.GetHeader("Accept")) {
		msg, _ := helper.GetTranslator(i, true).T("error.validation_error")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": msg, "errors": errors})
	} else {
		i.Session.Store.Flash("errors", errors)
		i.Redirector.Back(http.StatusSeeOther)
	}
}

func unauthenticated(c *gin.Context, err error, instance *instance.Instance) {
	if utils.ExpectJson(c.GetHeader("Accept")) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	} else {
		instance.Redirector.Guest(http.StatusSeeOther, "login")
	}
}

func handleTokenMismatch(c *gin.Context, err error) {
	if utils.ExpectJson(c.GetHeader("Accept")) {
		c.JSON(419, t.H{"message": err.Error()})
	} else {
		c.HTML(419, "errors/419.tmpl", t.Default(c).Wrap(t.H{"message": err.Error()}))
	}
}

func displayErrorWithStackTrace(c *gin.Context, code int, err stackTraceableError) {
	if utils.ExpectJson(c.GetHeader("Accept")) {
		c.JSON(code, convertStackTraceableError(err))
	} else if c.Writer.Written() {
		c.Data(code, "text/html; charset=utf-8", []byte(fmt.Sprintf("<p>%s</p>", err.Error())))
	} else {
		c.HTML(code, "errors/errors.tmpl", t.Default(c).Wrap(convertStackTraceableError(err)))
	}
}

func convertStackTraceableError(err stackTraceableError) map[string]any {
	trace := err.StackTrace()

	line, _ := strconv.Atoi(fmt.Sprintf("%d", trace[0]))

	return map[string]any{
		"message":  err.Error(),
		"function": fmt.Sprintf("%n", trace[0]),
		"file":     getFileFromStack(trace[0]),
		"line":     line,
		"trace":    formatStackTrace(trace),
	}
}

func formatStackTrace(traces []errors.Frame) []map[string]any {
	var lines = make([]map[string]any, 0, len(traces))
	for _, f := range traces {
		line, _ := strconv.Atoi(fmt.Sprintf("%d", f))

		lines = append(lines, map[string]any{
			"function": fmt.Sprintf("%n", f),
			"file":     getFileFromStack(f),
			"line":     line,
		})
	}
	return lines
}

func getFileFromStack(stack errors.Frame) string {
	return strings.SplitN(fmt.Sprintf("%+s", stack), "\n\t", 2)[1]
}

func displayErrorMessage(c *gin.Context, code int, message string) {
	if utils.ExpectJson(c.GetHeader("Accept")) {
		c.JSON(code, t.H{"message": message})
	} else if c.Writer.Written() {
		c.Data(code, "text/html; charset=utf-8", []byte(fmt.Sprintf("<p>%s</p>", message)))
	} else {
		c.HTML(code, "errors/errors.tmpl", t.Default(c).Wrap(t.H{"message": message}))
	}
}

type WriterMonitor struct {
	gin.ResponseWriter
	written bool
}

func (br *WriterMonitor) WriteString(s string) (int, error) {
	br.written = true
	return br.ResponseWriter.Write([]byte(s))
}

func (br *WriterMonitor) Write(data []byte) (int, error) {
	br.written = true
	return br.ResponseWriter.Write(data)
}

func (br *WriterMonitor) Written() bool {
	return br.written
}
