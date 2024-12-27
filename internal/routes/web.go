package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/auth/generic"
	authmid "github.com/wolftotem4/golava-core/auth/middleware"
	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/http/csrf"
	"github.com/wolftotem4/golava-core/instance"
	sessmid "github.com/wolftotem4/golava-core/session/middleware"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/middlewares"
	"github.com/wolftotem4/golava/internal/routes/home"
)

func LoadWebRoutes(r gin.IRouter, a *app.App) {
	r.Use(cookie.CookieMiddleware(a.Cookie))
	r.Use(sessmid.StartSession)
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

		r.GET("foo", func(c *gin.Context) {
			i := instance.MustGetInstance(c)
			err := i.Auth.(*generic.SessionGuard).LogoutOtherDevices(c, "test")
			if err != nil {
				c.Error(err)
				return
			}

			c.JSON(200, gin.H{
				"message": "foo",
			})
		})
	}
}
