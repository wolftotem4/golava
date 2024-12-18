package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava/internal/app"
)

func LoadApiRoutes(r gin.IRouter, app *app.App) {
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
