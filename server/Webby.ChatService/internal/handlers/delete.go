package handlers

import (
	"log/slog"
	"net/http"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

func (h *handler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "handlers.delete")

	var uri deleteUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.chatService.Delete(ctx, userID, chatID)
	if err != nil {
		log.Debug("chat deletion error", slog.String("err", err.Error()))
		HandleAppError(c, "failed to delete a chat", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Successfully deleted",
	})
}
