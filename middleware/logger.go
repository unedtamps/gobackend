package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := utils.NewULID()
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id.String())
		c.Next()
	}
}

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		reqID, _ := c.Get("request_id")

		log.Info("http_request",
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
		)
	}
}
