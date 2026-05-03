package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"webby-wsgateway/docs"
	"webby-wsgateway/internal/config"
	"webby-wsgateway/internal/ws"
	"webby-wsgateway/pkg/http/middleware/auth"
	"webby-wsgateway/pkg/http/middleware/cors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	loggerMw "webby-wsgateway/pkg/http/middleware/logger"
)

type Service interface {
	GenerateToken(ctx context.Context, userID uuid.UUID) (string, error)
}

func NewServer(cfg *config.Config, service Service, logger *slog.Logger, wsSrv *ws.Server) http.Handler {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	handler := NewHandler(service)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(loggerMw.Logger(logger), gin.Recovery(), cors.CORS())

	router.Any("/socket.io/*any", gin.WrapH(wsSrv.IO()))

	router.GET("/swagger", swaggerUI)
	router.GET("/swagger/", swaggerUI)
	router.GET("/swagger/index.html", swaggerUI)
	router.StaticFile("/swagger/doc.json", "./docs/oas.json")

	authGroup := router.Group("/")
	authGroup.Use(requireAuth)
	{
		authGroup.GET("/api/ws-token", handler.getWsToken)
	}

	return router
}

func swaggerUI(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docs.SwaggerUIHTML))
}
