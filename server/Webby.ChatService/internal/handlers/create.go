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
	RoomID *string `json:"roomId"`
}

type chatResponse struct {
	ID        uuid.UUID  `json:"id"`
	RoomID    *uuid.UUID `json:"roomId,omitempty"`
	CreatedAt string     `json:"createdAt"`
}

func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("op", "httpserver.chats.create")

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("invalid request body", slog.Any("error", err))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(*req.RoomID)
	chat, err := h.chatService.Create(ctx, &roomID)
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
			RoomID:    chat.RoomID,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
