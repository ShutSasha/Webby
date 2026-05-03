package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
	"webby-chat/internal/handlers/responses"
	"webby-chat/pkg/logger"

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
		responses.HandleValidationError(c, err)
		return
	}

	var roomId *uuid.UUID
	if req.RoomId != nil && strings.TrimSpace(*req.RoomId) != "" {
		parsed, err := uuid.Parse(strings.TrimSpace(*req.RoomId))
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors: map[string]string{"roomId": "invalid UUID format"},
			})
			return
		}
		roomId = &parsed
	}

	chat, err := h.service.Create(ctx, roomId)
	if err != nil {
		log.Error("create chat error", slog.String("err", err.Error()))
		responses.HandleAppError(c, "Create chat error", err)
		return
	}

	c.JSON(http.StatusCreated, responses.ApiResponse[chatResponse]{
		Success: true,
		Message: "Chat created",
		Data: &chatResponse{
			Id:        chat.Id,
			RoomId:    chat.RoomId,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
