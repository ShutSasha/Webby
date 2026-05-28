package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"webby/room-service/pkg/logger"

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
	log := logger.FromContext(ctx).With("operation", "handlers.ListPublic")

	var query listMyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("query validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	if strings.ToLower(query.Category) == "all" {
		query.Category = ""
	}

	log.Debug("Query", "page", query.Page, "limit", query.Limit, "search", query.Search, "category", query.Category)

	rooms, total, err := h.service.ListPublic(ctx, query.Page, query.Limit, query.Search, query.Category)
	if err != nil {
		log.Error("list public rooms error", slog.Any("err", err))
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
