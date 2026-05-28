package handlers

import (
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createRoomRequest struct {
	Name         string                `form:"name" binding:"required,min=2,max=50"`
	CategoryName string                `form:"categoryName" binding:"required"`
	IsPrivate    bool                  `form:"isPrivate"`
	Thumbnail    *multipart.FileHeader `form:"thumbnail"`
}

type createResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CategoryName string     `json:"categoryName"`
	IsPrivate    bool       `json:"isPrivate"`
	Thumbnail    string     `json:"thumbnail,omitempty"`
	ChatID       *uuid.UUID `json:"chatId,omitempty"`
}

func (h *handler) Create(c *gin.Context) {
	const maxFileSize = 2 * 1024 * 1024

	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With("operation", "handlers.Create")

	var req createRoomRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Debug("validation error", slog.Any("err", err))
		HandleValidationError(c, err)
		return
	}

	hostID, _ := uuid.Parse(ctx.Value("userID").(string))

	var thumbnailData []byte
	var thumbnailFilename string

	if req.Thumbnail != nil {
		if req.Thumbnail.Size > maxFileSize {
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
			log.Error("failed to open thumbnail", slog.String("error", err.Error()))
			HandleAppError(c, "File open error", err)
			return
		}
		defer file.Close()

		thumbnailData, err = io.ReadAll(file)
		if err != nil {
			log.Error("failed to read thumbnail", slog.String("error", err.Error()))
			HandleAppError(c, "File read error", err)
			return
		}
		thumbnailFilename = req.Thumbnail.Filename
	}

	room, err := h.service.Create(ctx, &models.Room{
		HostID:    hostID,
		Category:  req.CategoryName,
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
	}, thumbnailData, thumbnailFilename)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"categoryName": "category does not exist"},
			})
			return
		}
		log.Error("create room error", slog.Any("err", err))
		HandleAppError(c, "Create room error", err)
		return
	}

	c.JSON(http.StatusCreated, ApiResponse[createResponse]{
		Success: true,
		Message: "Room created",
		Data: &createResponse{
			ID:           room.ID,
			Name:         room.Name,
			CategoryName: room.Category,
			IsPrivate:    room.IsPrivate,
			Thumbnail:    room.Thumbnail,
			ChatID:       room.ChatID,
		},
	})
}
