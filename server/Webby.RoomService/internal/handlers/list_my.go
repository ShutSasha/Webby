package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listMyQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	Category string `form:"category" binding:"omitempty"`
}

type roomListMyItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"categoryName"`
	IsPrivate bool      `json:"isPrivate"`
	Thumbnail string    `json:"thumbnail"`
}

func (h *handler) ListMy(c *gin.Context) {
	ctx := c.Request.Context()

	var query listMyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleValidationError(c, err)
		return
	}

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	rooms, total, err := h.roomService.ListMyRooms(ctx, userID, query.Page, query.Limit, query.Search, query.Category)
	if err != nil {
		HandleAppError(c, "List rooms error", err)
		return
	}

	items := make([]roomListMyItem, len(rooms))
	for i, room := range rooms {
		items[i] = roomListMyItem{
			ID:        room.ID,
			Name:      room.Name,
			Category:  room.Category,
			IsPrivate: room.IsPrivate,
			Thumbnail: room.Thumbnail,
		}
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[roomListMyItem]]{
		Success: true,
		Message: "User rooms retrieved successfully",
		Data: &PaginatedResponse[roomListMyItem]{
			Items: items,
			Page:  query.Page,
			Limit: query.Limit,
			Total: int(total),
		},
	})
}
