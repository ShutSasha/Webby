package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type messageManager interface {
	SaveMessage(ctx context.Context, senderID, chatID uuid.UUID, text string) (*models.RichMessage, error)
}

type saveMessageUri struct {
	ChatID string `uri:"chatId" binding:"required,uuid"`
}

type saveMessageBody struct {
	SenderID uuid.UUID `json:"senderId" binding:"required"`
	Text     string    `json:"text" binding:"required"`
}

type sender struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarUrl"`
}

type saveMessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Sender    sender    `json:"sender"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *handler) SaveMessage(c *gin.Context) {
	const op = "handlers.SaveMessage"
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", op)

	var uri saveMessageUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	var body saveMessageBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Debug("body validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	message, err := h.messageManager.SaveMessage(ctx, body.SenderID, chatID, body.Text)
	if err != nil {
		log.Error("Error saving message", "err", err.Error())
		HandleAppError(c, "Error saving message", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[saveMessageResponse]{
		Success: true,
		Message: "Message saved",
		Data: &saveMessageResponse{
			ID: message.ID,
			Sender: sender{
				ID:        message.Sender.ID,
				Username:  message.Sender.Username,
				AvatarURL: message.Sender.AvatarURL,
			},
			Text:      message.Content,
			CreatedAt: message.CreatedAt,
		},
	})
}
