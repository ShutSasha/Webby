package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"webby/room-service/internal/apperrors"
	"webby/room-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type updateResponse struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryName string    `json:"categoryName"`
	IsPrivate    bool      `json:"isPrivate"`
	Thumbnail    string    `json:"thumbnail,omitempty"`
}

func (h *handler) Update(c *gin.Context) {
	const maxFileSize = 2 * 1024 * 1024

	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.update"))

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Debug("invalid UUID format", slog.String("id", idStr))
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  map[string]string{"id": "the id format is not valid"},
		})
		return
	}

	if err := c.Request.ParseMultipartForm(int64(maxFileSize)); err != nil {
		log.Debug("failed to parse multipart form", slog.String("error", err.Error()))
		HandleAppError(c, "Validation error", apperrors.ErrInvalidInput)
		return
	}

	var name, categoryName *string
	var isPrivate *bool
	problems := make(map[string]string)

	if nameVal := c.Request.FormValue("name"); nameVal != "" {
		if strings.TrimSpace(nameVal) == "" || len(nameVal) < 2 || len(nameVal) > 50 {
			problems["name"] = "must be between 2 and 50 characters and cannot be empty"
		} else {
			name = &nameVal
		}
	}

	if categoryVal := c.Request.FormValue("categoryName"); categoryVal != "" {
		if strings.TrimSpace(categoryVal) == "" {
			problems["categoryName"] = "category name cannot be whitespace-only"
		} else {
			categoryName = &categoryVal
		}
	}

	if isPrivateStr := c.Request.FormValue("isPrivate"); isPrivateStr != "" {
		isParsed, err := strconv.ParseBool(isPrivateStr)
		if err != nil {
			problems["isPrivate"] = "must be 'true' or 'false'"
		} else {
			isPrivate = &isParsed
		}
	}

	if len(problems) > 0 {
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  problems,
		})
		return
	}

	var thumbnailData *[]byte
	var thumbnailFilename *string

	file, fileHeader, err := c.Request.FormFile("thumbnail")
	if err == nil {
		defer file.Close()

		if fileHeader.Size > int64(maxFileSize) {
			c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"thumbnail": "file size must not exceed 2MB"},
			})
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			log.Error("failed to read thumbnail file", slog.String("error", err.Error()))
			HandleAppError(c, "File upload error", err)
			return
		}
		thumbnailData = &data
		filename := fileHeader.Filename
		thumbnailFilename = &filename
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.Debug("unexpected file error", slog.String("error", err.Error()))
		HandleAppError(c, "File upload error", apperrors.ErrInvalidInput)
		return
	}

	userIdStr := ctx.Value("userID").(string)
	userId, _ := uuid.Parse(userIdStr)

	updatedRoom, err := h.service.Update(ctx, id, name, categoryName, isPrivate, thumbnailData, thumbnailFilename, userId)
	if err != nil {
		log.Error("update room error", slog.Any("err", err))
		HandleAppError(c, "Update room error", err)
		return
	}

	c.JSON(http.StatusOK, ApiResponse[updateResponse]{
		Success: true,
		Message: "Room updated",
		Data: &updateResponse{
			Id:           updatedRoom.Id,
			Name:         updatedRoom.Name,
			CategoryName: updatedRoom.CategoryName,
			IsPrivate:    updatedRoom.IsPrivate,
			Thumbnail:    updatedRoom.Thumbnail,
		},
	})
}
