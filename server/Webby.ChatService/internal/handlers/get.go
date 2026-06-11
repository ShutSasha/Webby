package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chatGetUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

type chatGetResponse struct {
	ID        uuid.UUID  `json:"id"`
	RoomID    *uuid.UUID `json:"roomId,omitempty"`
	CreatedAt string     `json:"createdAt"`
}

func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	var uri chatGetUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	chat, err := h.chatService.GetByID(ctx, chatID)
	if err != nil {
		HandleAppError(c, "Get chat error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[chatGetResponse]{
		Success: true,
		Message: "Chat retrieved",
		Data: &chatGetResponse{
			ID:        chat.ID,
			RoomID:    chat.RoomID,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
