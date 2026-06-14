package handlers

import (
	"net/http"
	"time"
	"webby/chat-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chatGetUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

type chatGetResponse struct {
	ID        uuid.UUID   `json:"id"`
	CreatedAt string      `json:"createdAt"`
	User      models.User `json:"user"`
}

func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	var uri chatGetUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	chat, err := h.chatService.GetByID(ctx, chatID, userID)
	if err != nil {
		HandleAppError(c, "Get chat error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[chatGetResponse]{
		Success: true,
		Message: "Chat retrieved",
		Data: &chatGetResponse{
			ID:        chat.ID,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
			User:      chat.User,
		},
	})
}
