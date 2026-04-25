package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) ListMy(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.listMy"))

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

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	type roomListItem struct {
		Id           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		CategoryName string    `json:"categoryName"`
		IsPrivate    bool      `json:"isPrivate"`
		Thumbnail    string    `json:"thumbnail"`
	}

	rooms, total, err := h.service.ListMy(ctx, userId, page, limit, search, categoryName)
	if err != nil {
		log.Error("list my rooms error", slog.Any("err", err))
		HandleAppError(c, "List rooms error", err)
		return
	}

	items := make([]roomListItem, len(rooms))
	for i, room := range rooms {
		items[i] = roomListItem{
			Id:           room.Id,
			Name:         room.Name,
			CategoryName: room.CategoryName,
			IsPrivate:    room.IsPrivate,
			Thumbnail:    room.Thumbnail,
		}
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[roomListItem]]{
		Success: true,
		Message: "User rooms retrieved successfully",
		Data: &PaginatedResponse[roomListItem]{
			Items: items,
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	})
}
