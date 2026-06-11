package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type removeMemberUri struct {
	RoomID   string `uri:"id" binding:"required,uuid"`
	MemberID string `uri:"memberId" binding:"required,uuid"`
}

func (h *handler) RemoveMember(c *gin.Context) {
	ctx := c.Request.Context()

	var uri removeMemberUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	memberID, _ := uuid.Parse(uri.MemberID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if err := h.roomMemberService.RemoveMember(ctx, roomID, memberID, userID); err != nil {
		HandleAppError(c, "Remove member error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Member removed successfully",
	})
}
