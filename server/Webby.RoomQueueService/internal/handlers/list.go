package handlers

import (
	"net/http"

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

	var uri listUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var query listQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	items, total, err := h.service.GetQueue(ctx, roomID, userID, query.Page, query.Limit)
	if err != nil {
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
