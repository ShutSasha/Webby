package handlers

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateUri struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type updateRequest struct {
	Name      *string               `form:"name" binding:"omitempty,min=2,max=50"`
	Category  *string               `form:"categoryName" binding:"omitempty"`
	IsPrivate *bool                 `form:"isPrivate" binding:"omitempty"`
	Thumbnail *multipart.FileHeader `form:"thumbnail"`
}

type updateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"categoryName"`
	IsPrivate bool      `json:"isPrivate"`
	Thumbnail string    `json:"thumbnail,omitempty"`
}

func (h *handler) Update(c *gin.Context) {
	const maxFileSize = 2 * 1024 * 1024
	ctx := c.Request.Context()

	var uri updateUri
	if err := c.ShouldBindUri(&uri); err != nil {
		HandleValidationError(c, err)
		return
	}

	var req updateRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	roomID, _ := uuid.Parse(uri.RoomID)
	hostID, _ := uuid.Parse(ctx.Value("userID").(string))

	var thumbnailData *[]byte
	var thumbnailFilename *string

	if req.Thumbnail != nil {
		if req.Thumbnail.Size > int64(maxFileSize) {
			c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"thumbnail": "file size must not exceed 2MB"},
			})
			return
		}
		if req.Thumbnail.Size == 0 {
			c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"thumbnail": "file cannot be empty"},
			})
			return
		}

		file, err := req.Thumbnail.Open()
		if err != nil {
			HandleAppError(c, "File open error", err)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			HandleAppError(c, "File upload error", err)
			return
		}

		thumbnailData = &data
		filename := req.Thumbnail.Filename
		thumbnailFilename = &filename
	}

	updatedRoom, err := h.roomService.Update(
		ctx,
		roomID,
		hostID,
		req.Name,
		req.Category,
		thumbnailFilename,
		thumbnailData,
		req.IsPrivate,
	)
	if err != nil {
		HandleAppError(c, "Update room error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[updateResponse]{
		Success: true,
		Message: "Room updated",
		Data: &updateResponse{
			ID:        updatedRoom.ID,
			Name:      updatedRoom.Name,
			Category:  updatedRoom.Category,
			IsPrivate: updatedRoom.IsPrivate,
			Thumbnail: updatedRoom.Thumbnail,
		},
	})
}
