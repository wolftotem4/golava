package routes

import (
	"net/http"

	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/http/utils"
	"github.com/wolftotem4/golava-core/instance"
	langmid "github.com/wolftotem4/golava-core/lang/middleware"
	sessmid "github.com/wolftotem4/golava-core/session/middleware"
	t "github.com/wolftotem4/golava-core/template"
	tplmid "github.com/wolftotem4/golava-core/template/middleware"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/helper"
	"github.com/wolftotem4/golava/internal/middlewares"
	"golang.org/x/text/language"
)

func Register(r *gin.Engine, a *app.App) {
	// setup templates & functions
	tplmid.LoadFuncMap(r, a)
	r.LoadHTMLGlob("templates/**/*")

	// setup static files
	r.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))

	// setup global middlewares
	r.Use(
		brotli.Brotli(brotli.DefaultCompression),
		middlewares.Recovery(a.Debug),
		instance.NewInstance(a),
		sessmid.SaveSession,
		middlewares.ErrorHandle,
	)

	// setup host language support
	r.Use(langmid.SetLocale("hl", map[language.Tag]string{
		language.English: "en",
		// language.SimplifiedChinese:  "zh",
		// language.TraditionalChinese: "zh_Hant_TW",
	}))

	RegisterWebRoutes(r.Group("/"), a)
	// RegisterApiRoutes(r.Group("/api"), a)

	r.NoRoute(func(c *gin.Context) {
		if utils.ExpectJson(c.GetHeader("Accept")) {
			var i = instance.MustGetInstance(c)

			msg, _ := helper.GetTranslator(i, true).T("error.page_not_found")
			c.JSON(http.StatusNotFound, gin.H{"message": msg})
		} else {
			c.HTML(http.StatusNotFound, "errors/404.tmpl", t.Default(c))
		}
	})
}
