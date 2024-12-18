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
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/browser"
	"github.com/wolftotem4/golava-core/foundation"
	"github.com/wolftotem4/golava-core/instance"
	t "github.com/wolftotem4/golava-core/template"
	"github.com/wolftotem4/golava-core/validation"
	"github.com/wolftotem4/golava/internal/bootstrap"
	"github.com/wolftotem4/golava/middlewares"
	"github.com/wolftotem4/golava/routes"
)

func main() {
	ctx := context.Background()

	app, err := bootstrap.InitApp(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	binding.Validator = validation.NewMoldModValidator()

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
