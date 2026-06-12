package handlers

import (
	"net/http"
	"webby/vote-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var uri listUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	votings, err := h.service.ListVotings(ctx, roomID, userID)
	if err != nil {
		HandleAppError(c, "List votings error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[[]models.EnrichedVoting]{
		Success: true,
		Message: "Votes retrieved",
		Data:    &votings,
	})
}
