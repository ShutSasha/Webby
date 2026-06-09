package handlers

import (
	"log/slog"
	"net/http"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listQuery struct {
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
}

func (h *handler) history(c *gin.Context) {
	const op = "handlers.handler.list"
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", op)

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("query validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	chatHistory, total, err := h.chatService.History(ctx, userID, query.Page, query.Limit, query.Search)
	if err != nil {
		log.Debug("error retrieving user's chat history", slog.String("err", err.Error()))
		HandleAppError(c, "error retrieving user's chat history", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[models.ChatHistoryItem]]{
		Success: true,
		Message: "List chat history",
		Data: &PaginatedResponse[models.ChatHistoryItem]{
			Items: chatHistory,
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	})
}
