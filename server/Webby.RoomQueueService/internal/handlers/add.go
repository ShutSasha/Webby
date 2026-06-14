package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type addBody struct {
	VideoID string `json:"videoId" binding:"required,min=5"`
}

type addResponse struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
}

func (h *handler) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var uri addUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body addBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	itemID, position, err := h.service.AddToQueue(ctx, roomID, userID, body.VideoID)
	if err != nil {
		HandleAppError(c, "Add to queue error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[addResponse]{
		Success: true,
		Message: "Item added to queue",
		Data:    &addResponse{ID: itemID, Position: position},
	})
}
