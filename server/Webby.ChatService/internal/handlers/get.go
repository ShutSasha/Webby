package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"webby/chat-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chatGetResponse struct {
	Id        uuid.UUID  `json:"id"`
	RoomId    *uuid.UUID `json:"roomId,omitempty"`
	CreatedAt string     `json:"createdAt"`
}

func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.chats.get"),
	)

	chatIdStr := c.Param("id")
	chatId, err := uuid.Parse(chatIdStr)
	if err != nil {
		log.Debug("invalid chat id", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "invalid UUID format"},
		})
		return
	}

	chat, err := h.service.GetByID(ctx, chatId)
	if err != nil {
		log.Error("retrieve chat error", slog.String("err", err.Error()))
		HandleAppError(c, "Get chat error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[chatGetResponse]{
		Success: true,
		Message: "Chat retrieved",
		Data: &chatGetResponse{
			Id:        chat.ID,
			RoomId:    chat.RoomID,
			CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		},
	})
}
