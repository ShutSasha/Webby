package handlers

import (
	"net/http"
	"webby/admin-service/docs"
	"webby/admin-service/internal/config"
	"webby/admin-service/pkg/http/middleware/auth"

	"github.com/gin-gonic/gin"
)

func addRoutes(router *gin.Engine, cfg *config.Config, handler handler) {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	api := router.Group("/api")
	{
		protected := api.Group("")
		protected.Use(requireAuth)
		{
			protected.GET("/complaints", handler.List)
			protected.POST("/complaints/:id/accept", handler.Accept)
			protected.POST("/complaints/:id/deny", handler.Deny)
			protected.GET("/stats/registrations", handler.RegistrationsStats)
			protected.GET("/stats/subscriptions", handler.SubscriptionsStats)
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
