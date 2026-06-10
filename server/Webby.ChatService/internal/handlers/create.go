package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRequest struct {
	TargetID string `json:"targetId"`
}

type chatResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
}

func (h handler) create(c *gin.Context) {
	ctx := c.Request.Context()

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	initiatorID, _ := uuid.Parse(ctx.Value("userID").(string))
	targetID, _ := uuid.Parse(req.TargetID)
	chat, err := h.chatService.CreatePrivate(ctx, initiatorID, targetID)
	if err != nil {
		HandleAppError(c, "Create chat error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[chatResponse]{
		Success: true,
		Message: "Chat created",
		Data: &chatResponse{
			ID:        chat.ID,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
