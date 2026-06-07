package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"webby/room-service/pkg/logger"

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
	log := logger.FromContext(ctx).With("operation", "handlers.ListMy")

	var query listMyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("query validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	if strings.ToLower(query.Category) == "all" {
		query.Category = ""
	}

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	rooms, total, err := h.roomService.ListMyRooms(
		ctx,
		userID,
		query.Page, query.Limit, query.Search, query.Category,
	)
	if err != nil {
		log.Error("list my rooms error", slog.String("err", err.Error()))
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
