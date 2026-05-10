package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getResponse struct {
	Id           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CategoryName string     `json:"categoryName"`
	IsPrivate    bool       `json:"isPrivate"`
	HostId       uuid.UUID  `json:"hostId"`
	Thumbnail    string     `json:"thumbnail"`
	ChatId       *uuid.UUID `json:"chatId,omitempty"`
}

func (h *handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.get"))

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		log.Debug("invalid id", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "the id format is not valid"},
		})
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	room, err := h.service.GetById(ctx, roomId, userId)
	if err != nil {
		log.Error("retrieve room error", slog.String("err", err.Error()))
		HandleAppError(c, "Get room error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[getResponse]{
		Success: true,
		Message: "Room retrieved",
		Data: &getResponse{
			Id:           room.Id,
			Name:         room.Name,
			CategoryName: room.CategoryName,
			IsPrivate:    room.IsPrivate,
			HostId:       room.HostId,
			Thumbnail:    room.Thumbnail,
			ChatId:       room.ChatId,
		},
	})
}
