package handlers

import (
	"net/http"
	"webby/docs"
	"webby/internal/config"
	"webby/pkg/http/middleware/auth"

	"github.com/gin-gonic/gin"
)

func addRoutes(router *gin.Engine, cfg *config.Config, handler handler) {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	api := router.Group("/api")
	{
		api.GET("/categories", handler.List)

		protected := api.Group("/categories")
		protected.Use(requireAuth)
		{
			protected.POST("/", handler.Create)
			protected.PUT("/:name", handler.Update)
			protected.DELETE("/:name", handler.Delete)
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
