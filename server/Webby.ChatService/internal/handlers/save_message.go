package handlers

import (
	"log/slog"
	"net/http"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type saveMessageUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

type saveMessageBody struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (h *handler) saveMessage(c *gin.Context) {
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

	senderID, _ := uuid.Parse(ctx.Value("userID").(string))
	chatID, _ := uuid.Parse(uri.ChatID)
	err := h.messageService.SaveMessage(ctx, chatID, senderID, body.Content)
	if err != nil {
		log.Error("Error saving message", "err", err.Error())
		HandleAppError(c, "Error saving message", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Message sent",
	})
}
