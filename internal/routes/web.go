package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/foundation"
	"github.com/wolftotem4/golava-core/httputils/csrf"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/middlewares"
	"github.com/wolftotem4/golava/internal/routes/home"
)

func LoadWebRoutes(r gin.IRouter, a *app.App) {
	r.Use(cookie.CookieMiddleware(a.Cookie))
	r.Use(foundation.StartSession)
	r.Use(csrf.VerifyCsrfToken)
	r.Use(middlewares.WebAuth)

	r.GET("/", home.Homepage)

	// guest routes
	{
		r := r.Group("/")
		r.Use(foundation.RedirectIfAuthenticated("/"))

		r.GET("/login", home.Login)
		r.POST("/login", home.SubmitLogin)

		r.GET("/register", home.RegisterView)
		r.POST("/register", home.Register)
	}

	// protected routes
	{
		r := r.Group("/")
		r.Use(foundation.Authenticate)

		r.GET("/logout", home.SubmitLogout)
	}
}
