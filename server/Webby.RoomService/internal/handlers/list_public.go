package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listPublicQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	Category string `form:"category" binding:"omitempty"`
}

type roomListPublicItem struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	HostID        uuid.UUID `json:"hostId"`
	HostUsername  string    `json:"hostUsername"`
	HostAvatarUrl string    `json:"hostAvatarUrl"`
	Category      string    `json:"categoryName"`
	Thumbnail     string    `json:"thumbnail"`
}

func (h *handler) ListPublic(c *gin.Context) {
	ctx := c.Request.Context()

	var query listMyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleValidationError(c, err)
		return
	}

	rooms, total, err := h.roomService.ListPublicRooms(ctx, query.Page, query.Limit, query.Search, query.Category)
	if err != nil {
		HandleAppError(c, "List rooms error", err)
		return
	}

	items := make([]roomListPublicItem, len(rooms))
	for i, room := range rooms {
		items[i] = roomListPublicItem{
			ID:            room.ID,
			Name:          room.Name,
			HostID:        room.HostID,
			HostUsername:  room.HostUsername,
			HostAvatarUrl: room.HostAvatarUrl,
			Category:      room.Category,
			Thumbnail:     room.Thumbnail,
		}
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[roomListPublicItem]]{
		Success: true,
		Message: "Public rooms retrieved successfully",
		Data: &PaginatedResponse[roomListPublicItem]{
			Items: items,
			Page:  query.Page,
			Limit: query.Limit,
			Total: int(total),
		},
	})
}
