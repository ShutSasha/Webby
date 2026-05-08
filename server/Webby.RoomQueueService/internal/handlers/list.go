package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type queueItemResponse struct {
	ID        uuid.UUID `json:"id"`
	VideoID   string    `json:"videoId"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	VideoUrl  string    `json:"videoUrl"`
	IsActive  bool      `json:"isActive"`
	Position  int       `json:"position"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "handlers.List"),
	)

	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Debug("invalid room id", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors: map[string]string{
				"id": "the room id format is not valid",
			},
		})
		return
	}

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

	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	items, total, err := h.service.GetQueue(ctx, roomID, userID, page, limit)
	if err != nil {
		log.Error("list queue error", slog.Any("err", err))
		HandleAppError(c, "List queue error", err)
		return
	}

	result := make([]queueItemResponse, 0, len(items))
	for _, item := range items {
		entry := queueItemResponse{
			ID:        item.ID,
			VideoID:   item.VideoID,
			Title:     item.Title,
			Thumbnail: item.Thumbnail,
			VideoUrl:  item.VideoUrl,
			IsActive:  item.IsActive,
			Position:  item.Position,
		}
		result = append(result, entry)
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[queueItemResponse]]{
		Success: true,
		Message: "Queue retrieved successfully",
		Data: &PaginatedResponse[queueItemResponse]{
			Items: result,
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}
