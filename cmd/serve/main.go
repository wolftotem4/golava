package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"github.com/wolftotem4/golava-core/foundation"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/lang"
	t "github.com/wolftotem4/golava-core/template"
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

	foundation.LoadFuncMap(r, app)
	r.LoadHTMLGlob("templates/**/*")
	r.Use(brotli.Brotli(brotli.DefaultCompression))
	r.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))

	r.Use(instance.NewInstance(app))
	r.Use(foundation.SaveSession)
	r.Use(middlewares.ErrorHandle)
	r.Use(lang.SetLocale("hl", map[language.Tag]string{
		language.English:            "en",
		language.SimplifiedChinese:  "zh",
		language.TraditionalChinese: "zh_Hant_TW",
	}))

	routes.LoadWebRoutes(r, app)
	routes.LoadApiRoutes(r.Group("/api"), app)
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", t.Default(c))
	})

	appURL := os.Getenv("BASE_URL")
	app.Router.BaseURL, _ = app.Router.BaseURL.Parse(appURL)

	go func() {
		time.Sleep(1 * time.Second)
		browser.OpenURL(appURL)
	}()
	r.Run(os.Getenv("LISTEN_ADDR"))
}
