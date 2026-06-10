package handlers

import (
	"log/slog"
	"net/http"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listMessagesUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

type listMessagesQuery struct {
	Page  int `form:"page,default=1" binding:"omitempty,min=1"`
	Limit int `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

func (h *handler) listMessages(c *gin.Context) {
	const op = "handlers.handler.listMessages"
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", op)

	var uri listMessagesUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	var query listMessagesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("query validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	messages, total, err := h.messageService.List(ctx, chatID, userID, query.Limit, query.Page)
	if err != nil {
		log.Error("Failed to list messages", "err", err.Error())
		HandleAppError(c, "failed to list messages", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[models.RichMessage]]{
		Success: true,
		Message: "Messages successfully retrieved",
		Data: &PaginatedResponse[models.RichMessage]{
			Items: messages,
			Page:  query.Page,
			Limit: query.Limit,
			Total: int(total),
		},
	})
}
