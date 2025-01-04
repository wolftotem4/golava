package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
)

func LogRequests(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		timeTaken := time.Since(start)

		i := instance.MustGetInstance(c)
		logger.InfoContext(c, "request",
			slog.String("ip", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("useragent", c.Request.UserAgent()),
			slog.String("referer", c.Request.Referer()),
			slog.Int("status", c.Writer.Status()),
			slog.String("timetaken", timeTaken.String()),
			slog.Any("userid", i.Auth.ID()),
		)
	}
}
