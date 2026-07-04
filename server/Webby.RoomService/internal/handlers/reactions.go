package handlers

import (
	"net/http"
	"webby/room-service/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *handler) Reactions(c *gin.Context) {
	ctx := c.Request.Context()

	reactions, err := h.reactionService.List(ctx)
	if err != nil {
		HandleAppError(c, "Failed to list reactions", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[[]models.Reaction]{
		Success: true,
		Message: "Reactions retrieved",
		Data:    &reactions,
	})
}
