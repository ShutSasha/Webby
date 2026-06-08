package handlers

import (
	"net/http"
	"webby/chat-service/docs"
	"webby/chat-service/internal/config"
	"webby/chat-service/pkg/http/middleware/auth"

	"github.com/gin-gonic/gin"
)

func addRoutes(router *gin.Engine, cfg *config.Config, handler handler) {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	api := router.Group("/api")
	{
		chats := api.Group("/chats")
		chats.Use(requireAuth)
		{
			chats.POST("", handler.Create)
			chats.GET("/:id", handler.Get)
			chats.POST("/:id/messages", handler.saveMessage)
			chats.GET("/:id/messages", handler.listMessages)
			chats.PATCH("/:id/messages/:messageId", handler.updateMessage)
			chats.DELETE("/:id/messages/:messageId", handler.deleteMessage)
		}
	}

	router.GET("/swagger", swaggerUI)
	router.GET("/swagger/", swaggerUI)
	router.GET("/swagger/index.html", swaggerUI)
	router.GET("/swagger/doc.yaml", swaggerSpec)
}

func swaggerUI(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docs.SwaggerUIHTML))
}

func swaggerSpec(c *gin.Context) {
	c.Header("Content-Type", "application/x-yaml")
	c.File("./docs/oas.yml")
}
