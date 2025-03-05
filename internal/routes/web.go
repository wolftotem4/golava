package routes

import (
	"github.com/gin-gonic/gin"
	authmid "github.com/wolftotem4/golava-core/auth/middleware"
	"github.com/wolftotem4/golava-core/http/csrf"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/middlewares"
	"github.com/wolftotem4/golava/internal/routes/home"
)

func RegisterWebRoutes(r gin.IRouter, a *app.App) {
	r.Use(middlewares.Cookie)
	r.Use(middlewares.StartSession)
	r.Use(csrf.VerifyCsrfToken)
	r.Use(middlewares.WebAuth)
	// r.Use(authmid.AuthenticateSession)
	r.Use(func(c *gin.Context) {
		i := instance.MustGetInstance(c)
		if i.Auth.Check() {
			c.Header("Cache-Control", "private")
		}
	})

	r.GET("/", home.Homepage)

	// guest routes
	{
		r := r.Group("/")
		r.Use(authmid.RedirectIfAuthenticated("/"))

		r.GET("/login", home.Login)
		r.POST("/login", home.SubmitLogin)

		r.GET("/register", home.RegisterView)
		r.POST("/register", home.Register)
	}

	// protected routes
	{
		r := r.Group("/")
		r.Use(authmid.Authenticate)

		r.GET("/logout", home.SubmitLogout)
	}
}
