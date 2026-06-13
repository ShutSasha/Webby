package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type voteForNextUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type voteForNextBody struct {
	QueueItemID uuid.UUID `uri:"queueItemId" binding:"required"`
}

func (h *handler) VoteForNext(c *gin.Context) {
	ctx := c.Request.Context()

	var uri voteForNextUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body voteForNextBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.service.VoteForNextVideo(ctx, roomID, userID, body.QueueItemID)
	if err != nil {
		HandleAppError(c, "failed to create voting for next video", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Successfully voted",
	})
}
