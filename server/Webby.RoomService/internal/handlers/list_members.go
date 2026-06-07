package handlers

import (
	"log/slog"
	"net/http"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listMembersUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type listMembersQuery struct {
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
}

type memberItem struct {
	UserID     uuid.UUID `json:"userId"`
	Username   string    `json:"username"`
	AvatarUrl  string    `json:"avatarUrl"`
	RoomPoints int       `json:"roomPoints"`
}

func (h *handler) ListMembers(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.ListMembers")

	var uri listMembersUri
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Debug("uri validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	var query listMembersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Debug("query validation error", slog.String("err", err.Error()))
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	members, total, err := h.roomMemberService.ListMembers(ctx, roomID, query.Page, query.Limit, query.Search)
	if err != nil {
		log.Error("list members error", slog.String("err", err.Error()))
		HandleAppError(c, "List room members error", err)
		return
	}

	items := make([]memberItem, len(members))
	for i, member := range members {
		items[i] = memberItem{
			UserID:     member.UserID,
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
			Page:  query.Page,
			Limit: query.Limit,
			Total: int(total),
		},
	})
}
