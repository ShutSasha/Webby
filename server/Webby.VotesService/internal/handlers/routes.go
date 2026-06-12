package handlers

import (
	"net/http"
	"webby/vote-service/docs"
	"webby/vote-service/internal/config"
	"webby/vote-service/pkg/http/middleware/auth"

	"github.com/gin-gonic/gin"
)

func addRoutes(router *gin.Engine, cfg *config.Config, handler handler) {
	requireAuth := auth.AuthMiddleware([]byte(cfg.JwtSecret))

	api := router.Group("/api/rooms/:id/votes")
	api.Use(requireAuth)
	{
		api.GET("", handler.List)
		api.POST("/right-choice", handler.CreateRightChoiceVoting)
		api.POST("/:voteId/resolve", handler.ResolveVoting)
		// api.POST("/next-video", handler.CreateNextVideoVoting)
		api.POST("/:voteId/vote", handler.CastVote)
	}

	router.GET("/swagger", swaggerUI)
	router.GET("/swagger/", swaggerUI)
	router.GET("/swagger/index.html", swaggerUI)
	router.GET("/swagger/doc.yaml", swaggerSpec)
}

func swaggerUI(c *gin.Context) {
	c.Data(
		http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(docs.SwaggerUIHTML),
	)
}

func swaggerSpec(c *gin.Context) {
	c.Header("Content-Type", "application/x-yaml")
	c.File("./docs/oas.yml")
}
