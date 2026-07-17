package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type useReactionUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type useReactionBody struct {
	ReactionID uuid.UUID `json:"reactionId" binding:"required,uuid"`
}

func (h *handler) UseReaction(c *gin.Context) {
	ctx := c.Request.Context()

	var uri useReactionUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body useReactionBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.reactionService.UseReaction(ctx, roomID, userID, body.ReactionID)
	if err != nil {
		HandleAppError(c, "Failed to use reaction", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Reaction sent",
	})
}
