package handlers

import (
	"log/slog"
	"net/http"

	"webby/wsgateway/docs"
	"webby/wsgateway/internal/config"
	"webby/wsgateway/internal/ws"
	"webby/wsgateway/pkg/http/middleware/auth"
	"webby/wsgateway/pkg/http/middleware/cors"

	"github.com/gin-gonic/gin"

	loggerMw "webby/wsgateway/pkg/http/middleware/logger"
)

func NewServer(cfg *config.Config, service tokenGenerator, logger *slog.Logger, wsSrv *ws.Server, broker broker) http.Handler {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	handler := New(service, broker)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(loggerMw.Logger(logger), gin.Recovery(), cors.CORS())

	router.Any("/socket.io/*any", gin.WrapH(wsSrv.IO()))

	router.GET("/swagger", swaggerUI)
	router.GET("/swagger/", swaggerUI)
	router.GET("/swagger/index.html", swaggerUI)
	router.StaticFile("/swagger/doc.json", "./docs/oas.yaml")

	authGroup := router.Group("/")
	authGroup.Use(requireAuth)
	{
		authGroup.GET("/api/ws-token", handler.getWsToken)
		authGroup.GET("/api/sse", handler.streamSSE)
	}

	return router
}

func swaggerUI(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docs.SwaggerUIHTML))
}
