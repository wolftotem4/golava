package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/lang"
	sessmid "github.com/wolftotem4/golava-core/session/middleware"
	t "github.com/wolftotem4/golava-core/template"
	tplmid "github.com/wolftotem4/golava-core/template/middleware"
	"github.com/wolftotem4/golava/internal/bootstrap"
	"github.com/wolftotem4/golava/internal/middlewares"
	"github.com/wolftotem4/golava/internal/routes"
	"golang.org/x/text/language"
)

func main() {
	ctx := context.Background()

	app, err := bootstrap.InitApp(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	gin.SetMode(os.Getenv(gin.EnvGinMode))
	r := gin.New()
	r.Use(gin.Recovery())

	tplmid.LoadFuncMap(r, app)
	r.LoadHTMLGlob("templates/**/*")
	r.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))

	r.Use(
		instance.NewInstance(app),
		brotli.Brotli(brotli.DefaultCompression),
		sessmid.SaveSession,
		middlewares.ErrorHandle,
	)

	r.Use(lang.SetLocale("hl", map[language.Tag]string{
		language.English:            "en",
		language.SimplifiedChinese:  "zh",
		language.TraditionalChinese: "zh_Hant_TW",
	}))

	routes.LoadWebRoutes(r.Group("/"), app)
	// routes.LoadApiRoutes(r.Group("/api"), app)

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", t.Default(c))
	})

	appURL := os.Getenv("BASE_URL")
	app.Router.BaseURL, _ = app.Router.BaseURL.Parse(appURL)

	r.Run(os.Getenv("LISTEN_ADDR"))
}
