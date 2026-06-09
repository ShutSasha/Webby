package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"webby/chat-service/pkg/logger"

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
	log := logger.FromContext(ctx).With("op", "handlers.handler.create")

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("invalid request body", slog.Any("error", err))
		HandleValidationError(c, err)
		return
	}

	initiatorID, _ := uuid.Parse(ctx.Value("userID").(string))
	targetID, _ := uuid.Parse(req.TargetID)
	chat, err := h.chatService.CreatePrivate(ctx, initiatorID, targetID)
	if err != nil {
		log.Error("create chat error", slog.String("err", err.Error()))
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
