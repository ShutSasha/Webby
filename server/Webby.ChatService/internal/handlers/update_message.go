package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateMessageUri struct {
	ChatID    string `uri:"id" binding:"required,uuid"`
	MessageID string `uri:"messageId" binding:"required,uuid"`
}

type updateMessageBody struct {
	Content string `json:"content" binding:"required,required,min=1"`
}

func (h *handler) updateMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var uri updateMessageUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var body updateMessageBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleValidationError(c, err)
		return
	}

	chatID, _ := uuid.Parse(uri.ChatID)
	messageID, _ := uuid.Parse(uri.MessageID)
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	err := h.messageService.Update(ctx, chatID, messageID, userID, body.Content)
	if err != nil {
		HandleAppError(c, "failed to update messages", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[struct{}]{
		Success: true,
		Message: "Message updated",
	})
}
