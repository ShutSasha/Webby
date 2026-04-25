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
		api.GET("/rooms/public", handler.ListPublic)

		rooms := api.Group("/rooms")
		rooms.Use(requireAuth)
		{
			rooms.POST("", handler.Create)
			rooms.GET("/my", handler.ListMy)
			rooms.GET("/:id", handler.Get)
			rooms.PUT("/:id", handler.Update)
			rooms.DELETE("/:id", handler.Delete)
			rooms.GET("/:id/members", handler.ListMembers)
			rooms.POST("/:id/members", handler.AddMembers)
			rooms.DELETE("/:id/members/:memberId", handler.RemoveMember)
			rooms.PATCH("/:id/members/:memberId/points", handler.UpdatePoints)
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
