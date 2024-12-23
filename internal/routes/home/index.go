package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	t "github.com/wolftotem4/golava-core/template"
)

func Homepage(c *gin.Context) {
	c.HTML(http.StatusOK, "home/index.tmpl", t.Default(c))
}
