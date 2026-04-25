package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"webby/internal/apperrors"
	"webby/internal/models"
	"webby/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createResponse struct {
	Id           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	CategoryName string     `json:"categoryName"`
	IsPrivate    bool       `json:"isPrivate"`
	Thumbnail    string     `json:"thumbnail,omitempty"`
	ChatId       *uuid.UUID `json:"chatId,omitempty"`
}

func (h *handler) Create(c *gin.Context) {
	const maxFileSize = 2 * 1024 * 1024

	ctx := c.Request.Context()
	log := logger.FromContext(ctx).With(slog.String("operation", "httpserver.rooms.create"))

	if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
		log.Debug("failed to parse multipart form", slog.String("error", err.Error()))
		HandleAppError(c, "Validation error", apperrors.ErrInvalidInput)
		return
	}

	name := c.Request.FormValue("name")
	categoryName := c.Request.FormValue("categoryName")
	isPrivateStr := c.Request.FormValue("isPrivate")

	problems := make(map[string]string)

	if strings.TrimSpace(name) == "" || len(name) < 2 || len(name) > 50 {
		problems["name"] = "must be between 2 and 50 characters and cannot be empty"
	}

	if strings.TrimSpace(categoryName) == "" {
		problems["categoryName"] = "category name is required"
	}

	isPrivate, err := strconv.ParseBool(isPrivateStr)
	if err != nil {
		problems["isPrivate"] = "must be 'true' or 'false'"
	}

	if len(problems) > 0 {
		c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
			Success: false,
			Message: "Validation error",
			Errors:  problems,
		})
		return
	}

	hostIdStr := ctx.Value("userID").(string)
	hostId, _ := uuid.Parse(hostIdStr)

	var thumbnailData []byte
	var thumbnailFilename string

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

		if fileHeader.Size == 0 {
			c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
				Success: false,
				Message: "Validation error",
				Errors:  map[string]string{"thumbnail": "file cannot be empty"},
			})
			return
		}

		thumbnailData, err = io.ReadAll(file)
		if err != nil {
			log.Error("failed to read thumbnail file", slog.String("error", err.Error()))
			HandleAppError(c, "File upload error", err)
			return
		}
		thumbnailFilename = fileHeader.Filename
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.Debug("unexpected file error", slog.String("error", err.Error()))
		HandleAppError(c, "File upload error", apperrors.ErrInvalidInput)
		return
	}

	room, err := h.service.Create(ctx, &models.Room{
		HostId:       hostId,
		CategoryName: categoryName,
		Name:         name,
		IsPrivate:    isPrivate,
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
			Id:           room.Id,
			Name:         room.Name,
			CategoryName: room.CategoryName,
			IsPrivate:    room.IsPrivate,
			Thumbnail:    room.Thumbnail,
			ChatId:       room.ChatId,
		},
	})
}
