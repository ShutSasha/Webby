package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type getResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CategoryName string     `json:"categoryName"`
	IsPrivate    bool       `json:"isPrivate"`
	HostID       uuid.UUID  `json:"hostId"`
	Thumbnail    string     `json:"thumbnail"`
	ChatID       *uuid.UUID `json:"chatId,omitempty"`
}

func (h *handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	var uri getUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	room, err := h.roomService.AccessRoom(ctx, roomID, userID)
	if err != nil {
		HandleAppError(c, "Get room error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[getResponse]{
		Success: true,
		Message: "Room retrieved",
		Data: &getResponse{
			ID:           room.ID,
			Name:         room.Name,
			CategoryName: room.Category,
			IsPrivate:    room.IsPrivate,
			HostID:       room.HostID,
			Thumbnail:    room.Thumbnail,
			ChatID:       room.ChatID,
		},
	})
}
