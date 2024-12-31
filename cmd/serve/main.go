package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wolftotem4/golava-core/instance"
	langmid "github.com/wolftotem4/golava-core/lang/middleware"
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

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Initialize logger
	err = bootstrap.InitLogger()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Initialize app
	app, err := bootstrap.InitApp(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Initialize gin
	gin.SetMode(os.Getenv(gin.EnvGinMode))
	r := gin.New()
	r.Use(gin.Recovery())

	// setup templates & functions
	tplmid.LoadFuncMap(r, app)
	r.LoadHTMLGlob("templates/**/*")

	// setup static files
	r.Use(static.Serve("/assets", static.LocalFile("./public/assets", true)))

	// setup global middlewares
	r.Use(
		instance.NewInstance(app),
		brotli.Brotli(brotli.DefaultCompression),
		sessmid.SaveSession,
		middlewares.ErrorHandle,
	)

	// setup host language support
	r.Use(langmid.SetLocale("hl", map[language.Tag]string{
		language.English: "en",
		// language.SimplifiedChinese:  "zh",
		// language.TraditionalChinese: "zh_Hant_TW",
	}))

	routes.LoadWebRoutes(r.Group("/"), app)
	// routes.LoadApiRoutes(r.Group("/api"), app)

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", t.Default(c))
	})

	appURL := os.Getenv("BASE_URL")
	app.Router.BaseURL, _ = url.Parse(appURL)

	r.Run(os.Getenv("LISTEN_ADDR"))
}
