package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getMemberPointsUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type getMemberPointsResponse struct {
	Points int `json:"points"`
}

func (h *handler) GetMemberPoints(c *gin.Context) {
	ctx := c.Request.Context()

	var uri getMemberPointsUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	points, err := h.roomMemberService.GetMemberPoints(ctx, roomID, userID)
	if err != nil {
		HandleAppError(c, "failed to get member points", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[getMemberPointsResponse]{
		Success: true,
		Message: "Successfully retrieve member points",
		Data:    &getMemberPointsResponse{Points: points},
	})
}
