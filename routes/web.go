package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/foundation"
	"github.com/wolftotem4/golava-core/httputils/csrf"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/middlewares"
	"github.com/wolftotem4/golava/routes/home"
	"github.com/wolftotem4/golava/routes/user"
)

func LoadWebRoutes(r gin.IRouter, app *app.App) {
	r.Use(cookie.CookieMiddleware(app.Cookie))
	r.Use(foundation.StartSession)
	r.Use(csrf.VerifyCsrfToken)
	r.Use(middlewares.WebAuth)

	r.GET("/", home.Homepage)

	{
		r := r.Group("/")
		r.Use(foundation.RedirectIfAuthenticated("/"))

		r.GET("/login", home.Login)
		r.POST("/login", home.SubmitLogin)

		r.GET("/register", user.RegisterView)
		r.POST("/register", user.Register)
	}

	{
		r := r.Group("/")
		r.Use(foundation.Authenticate)

		r.GET("/logout", home.SubmitLogout)
	}
}
