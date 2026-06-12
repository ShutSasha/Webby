package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type resolveVotingUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
	VoteID string `uri:"voteId" binding:"required,uuid"`
}

type resolveVotingRequest struct {
	RightChoice string `json:"rightChoice" binding:"required,min=1,max=100"`
}

func (h *handler) ResolveVoting(c *gin.Context) {
	ctx := c.Request.Context()

	var uri resolveVotingUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var req resolveVotingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	voteID, _ := uuid.Parse(uri.VoteID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.service.ResolveVoting(ctx, roomID, userID, voteID, req.RightChoice)
	if err != nil {
		HandleAppError(c, "Resolve vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Voting resolved and results published",
	})
}
