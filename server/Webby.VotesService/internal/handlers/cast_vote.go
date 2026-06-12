package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type castVoteUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
	VoteID string `uri:"voteId" binding:"required,uuid"`
}

type castVoteRequest struct {
	Choice string `json:"choice" binding:"required,min=1,max=100"`
}

func (h *handler) CastVote(c *gin.Context) {
	ctx := c.Request.Context()

	var uri castVoteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var req castVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	voteID, _ := uuid.Parse(uri.VoteID)
	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.service.CastVote(ctx, roomID, voteID, userID, req.Choice)
	if err != nil {
		HandleAppError(c, "Cast vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Vote accepted",
	})
}
