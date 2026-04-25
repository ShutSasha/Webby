package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handler) ListMembers(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.listMembers"))

	roomIdStr := c.Param("id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		log.Debug("invalid id", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "the id format is not valid"},
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

	search := c.Query("search")

	type memberItem struct {
		UserId     uuid.UUID `json:"userId"`
		Username   string    `json:"username"`
		AvatarUrl  string    `json:"avatarUrl"`
		RoomPoints int       `json:"roomPoints"`
	}

	members, total, err := h.service.ListMembers(ctx, roomId, page, limit, search)
	if err != nil {
		log.Error("list members error", slog.Any("err", err))
		HandleAppError(c, "List room members error", err)
		return
	}

	items := make([]memberItem, len(members))
	for i, member := range members {
		items[i] = memberItem{
			UserId:     member.UserId,
			Username:   member.Username,
			AvatarUrl:  member.AvatarUrl,
			RoomPoints: member.RoomPoints,
		}
	}

	c.JSON(http.StatusOK, ApiResponse[PaginatedResponse[memberItem]]{
		Success: true,
		Message: "Room members retrieved successfully",
		Data: &PaginatedResponse[memberItem]{
			Items: items,
			Page:  page,
			Limit: limit,
			Total: int(total),
		},
	})
}
