package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type saveMessageUri struct {
	ChatID string `uri:"id" binding:"required,uuid"`
}

type saveMessageBody struct {
	Content string `json:"content" binding:"required,min=1"`
}

func (h *handler) saveMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var uri saveMessageUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body saveMessageBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	senderID, _ := uuid.Parse(ctx.Value("userID").(string))
	chatID, _ := uuid.Parse(uri.ChatID)
	err := h.messageService.SaveMessage(ctx, chatID, senderID, body.Content)
	if err != nil {
		HandleAppError(c, "Error saving message", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Message sent",
	})
}
