package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) ListPublic(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.listPublic"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	search := c.Query("search")

	var categoryName *string
	if categoryStr := c.Query("category"); categoryStr != "" {
		if strings.ToLower(categoryStr) != "all" {
			categoryName = &categoryStr
		}
	}

	type roomListItem struct {
		Id            uuid.UUID `json:"id"`
		Name          string    `json:"name"`
		HostId        uuid.UUID `json:"hostId"`
		HostUsername  string    `json:"hostUsername"`
		HostAvatarUrl string    `json:"hostAvatarUrl"`
		CategoryName  string    `json:"categoryName"`
		Thumbnail     string    `json:"thumbnail"`
	}

	rooms, total, err := h.service.ListPublic(ctx, page, limit, search, categoryName)
	if err != nil {
		log.Error("list public rooms error", slog.Any("err", err))
		HandleAppError(c, "List rooms error", err)
		return
	}

	items := make([]roomListItem, len(rooms))
	for i, room := range rooms {
		items[i] = roomListItem{
			Id:            room.Id,
			Name:          room.Name,
			HostId:        room.HostId,
			HostUsername:  room.HostUsername,
			HostAvatarUrl: room.HostAvatarUrl,
			CategoryName:  room.CategoryName,
			Thumbnail:     room.Thumbnail,
		}
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[roomListItem]]{
		Success: true,
		Message: "Public rooms retrieved successfully",
		Data: &PaginatedResponse[roomListItem]{
			Items: items,
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	})
}
