package handlers

import (
	"net/http"
	"webby/vote-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type hasNextVideoVotingUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) HasNextVideoVoting(c *gin.Context) {
	ctx := c.Request.Context()

	var uri hasNextVideoVotingUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	voting, err := h.service.HasNextVideoVoting(ctx, roomID, userID)
	if err != nil {
		HandleAppError(c, "failed to check for next video voting", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[models.NextVideoInfo]{
		Success: true,
		Message: "Checked if room has next video voting",
		Data:    voting,
	})
}
