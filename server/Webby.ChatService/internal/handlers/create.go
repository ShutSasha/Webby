package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"webby/chat-service/internal/models"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRequest struct {
	RoomId *string `json:"roomId"`
}

type chatResponse struct {
	Id        uuid.UUID  `json:"id"`
	RoomId    *uuid.UUID `json:"roomId,omitempty"`
	CreatedAt string     `json:"createdAt"`
}

func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.chats.create"),
	)

	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("invalid request body", slog.Any("error", err))
		HandleValidationError(c, err)
		return
	}

	chat, err := h.service.Create(ctx, models.CreateChatRequest(req))
	if err != nil {
		log.Error("create chat error", slog.String("err", err.Error()))
		HandleAppError(c, "Create chat error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[chatResponse]{
		Success: true,
		Message: "Chat created",
		Data: &chatResponse{
			Id:        chat.Id,
			RoomId:    chat.RoomId,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
