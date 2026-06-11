package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRightChoiceVotingUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type createRightChoiceVotingRequest struct {
	VoteText string   `json:"voteText" binding:"required,min=1,max=500"`
	Duration int      `json:"duration" binding:"required,min=30,max=7200"`
	Choices  []string `json:"choices" binding:"required,min=2,max=10,dive"`
}

func (h *handler) CreateRightChoiceVoting(c *gin.Context) {
	ctx := c.Request.Context()

	var uri createRightChoiceVotingUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var req createRightChoiceVotingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.service.CreateWithRightChoice(ctx, roomID, userID, req.VoteText, req.Duration, req.Choices)
	if err != nil {
		HandleAppError(c, "Create vote error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Vote created",
	})
}
