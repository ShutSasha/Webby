package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type listQuery struct {
	Page  int `form:"page,default=1" binding:"omitempty,min=1"`
	Limit int `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

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
	log := logger.FromContext(ctx).With("operation", "handlers.List")

	var uri listUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("invalid room id in uri", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "invalid room id format"},
		})
		return
	}

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("invalid query params", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"query": "invalid page or limit values"},
		})
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	items, total, err := h.service.GetQueue(ctx, roomID, userID, query.Page, query.Limit)
	if err != nil {
		log.Error("list queue error", slog.String("err", err.Error()))
		HandleAppError(c, "List queue error", err)
		return
	}

	result := make([]queueItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, queueItemResponse{
			ID:        item.ID,
			VideoID:   item.VideoID,
			Title:     item.Title,
			Thumbnail: item.Thumbnail,
			VideoUrl:  item.VideoUrl,
			IsActive:  item.IsActive,
			Position:  item.Position,
		})
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[queueItemResponse]]{
		Success: true,
		Message: "Queue retrieved successfully",
		Data: &PaginatedResponse[queueItemResponse]{
			Items: result,
			Page:  query.Page,
			Limit: query.Limit,
			Total: total,
		},
	})
}
