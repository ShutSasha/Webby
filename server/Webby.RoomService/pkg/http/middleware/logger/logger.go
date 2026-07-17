package logger

import (
	"log/slog"

	ctxlogger "webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")

		log.Info(
			"request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("request_id", reqID),
		)

		ctx := ctxlogger.ToContext(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
