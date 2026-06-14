package handlers

import (
	"log/slog"
	"net/http"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteMessageUri struct {
	ChatID    string `uri:"id" binding:"required,uuid"`
	MessageID string `uri:"messageId" binding:"required,uuid"`
}

func (h *handler) deleteMessage(c *gin.Context) {
	const op = "handlers.hanler.deleteMessage"
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", op)

	var uri deleteMessageUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	messageID, _ := uuid.Parse(uri.MessageID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.messageService.Delete(ctx, chatID, messageID, userID)
	if err != nil {
		log.Error("failed to delete messages", "err", err.Error())
		HandleAppError(c, "failed to delete messages", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Message deleted",
	})
}
