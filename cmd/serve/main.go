package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/wolftotem4/golava/internal/bootstrap"
	"github.com/wolftotem4/golava/internal/env"
	"github.com/wolftotem4/golava/internal/routes"
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
	err = bootstrap.InitDefaultLogger()
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

	// Setup loggers
	err = bootstrap.InitLogger(app,
		bootstrap.Logger("app", slog.Default()),
		bootstrap.LoggerSink("request", env.Get("REQUEST_LOG_SINK", "console")),
	)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	// Initialize gin
	gin.SetMode(os.Getenv(gin.EnvGinMode))
	r := gin.New()
	r.Use(gin.Recovery())

	routes.Register(r, app)

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.InfoContext(ctx, fmt.Sprintf("Listening on %s...", listenAddr))
	if err := r.Run(listenAddr); err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
}
