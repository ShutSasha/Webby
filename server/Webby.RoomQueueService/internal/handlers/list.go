package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"webby-room-queue/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type videoChild struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	VideoUrl  string    `json:"videoUrl"`
}

type queueItemResponse struct {
	Id            uuid.UUID    `json:"id"`
	EntityId      uuid.UUID    `json:"entityId"`
	EntityType    string       `json:"entityType"`
	Title         string       `json:"title"`
	Thumbnail     string       `json:"thumbnail"`
	VideoUrl      string       `json:"videoUrl"`
	IsActive      bool         `json:"isActive"`
	IsFolder      bool         `json:"isFolder"`
	Position      int          `json:"position"`
	TotalChildren int          `json:"totalChildren"`
	Children      []videoChild `json:"children,omitempty"`
}

func (h *handler) List(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(
		slog.String("operation", "httpserver.queue.list"),
	)

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
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

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	items, total, err := h.service.GetQueue(
		ctx, roomId, userId, page, limit,
	)
	if err != nil {
		log.Error("list queue error", slog.Any("err", err))
		HandleAppError(c, "List queue error", err)
		return
	}

	result := make([]queueItemResponse, 0, len(items))
	for _, item := range items {
		entry := queueItemResponse{
			Id:            item.Id,
			EntityId:      item.EntityId,
			EntityType:    item.EntityType,
			Title:         item.Title,
			Thumbnail:     item.Thumbnail,
			VideoUrl:      item.VideoUrl,
			IsActive:      item.IsActive,
			IsFolder:      item.IsFolder,
			Position:      item.Position,
			TotalChildren: item.TotalChildren,
		}
		if len(item.Children) > 0 {
			children := make([]videoChild, 0, len(item.Children))
			for _, ch := range item.Children {
				children = append(children, videoChild{
					Id:        ch.Id,
					Title:     ch.Title,
					Thumbnail: ch.Thumbnail,
					VideoUrl:  ch.VideoUrl,
				})
			}
			entry.Children = children
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
