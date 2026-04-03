package update

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/models"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Updater interface {
	Update(ctx context.Context, roomId uuid.UUID, name *string, categoryName *string, isPrivate *bool, thumbnailData []byte, thumbnailFilename string, userId uuid.UUID) (*models.Room, error)
}

func New(logger *slog.Logger, updater Updater) http.Handler {
	return errorWrapper.MakeHandler(logger, updateRoom(logger, updater))
}

// @Title Update a room
// @Description Update a room's properties. You can update any combination of fields (name, categoryName, isPrivate, thumbnail). Only the room creator can update the room. Fields not provided in the request will not be changed.
// @Param  id          path  string  true   "Room ID (UUID v4 format)"
// @Param  name        form  string  false  "Room name (2-50 chars, optional)"
// @Param  categoryName form  string  false  "Category name (optional)"
// @Param  isPrivate   form  string  false  "Visibility flag ('true' or 'false', optional)"
// @Param  thumbnail   file  file    false  "Room thumbnail image (optional, max 2MB)"
// @Success  200  object docs.RoomApiResponse  "Room successfully updated"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized"
// @Failure  404  object docs.ErrorResponse  "Room not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id} [put]
func updateRoom(logger *slog.Logger, updater Updater) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.update"))

	const maxFileSize = 2 * 1024 * 1024 // 2 MB

	type response struct {
		Id           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		CategoryName string    `json:"categoryName"`
		IsPrivate    bool      `json:"isPrivate"`
		Thumbnail    string    `json:"thumbnail,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Debug("invalid UUID format", slog.String("id", idStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			log.Debug("failed to parse multipart form", slog.String("error", err.Error()))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		var name, categoryName *string
		var isPrivate *bool
		problems := make(map[string]string)

		if nameVal := r.FormValue("name"); nameVal != "" {
			if strings.TrimSpace(nameVal) == "" || len(nameVal) < 2 || len(nameVal) > 50 {
				problems["name"] = "must be between 2 and 50 characters and cannot be empty"
			} else {
				name = &nameVal
			}
		}

		if categoryVal := r.FormValue("categoryName"); categoryVal != "" {
			if strings.TrimSpace(categoryVal) == "" {
				problems["categoryName"] = "category name cannot be whitespace-only"
			} else {
				categoryName = &categoryVal
			}
		}

		if isPrivateStr := r.FormValue("isPrivate"); isPrivateStr != "" {
			isParsed, err := strconv.ParseBool(isPrivateStr)
			if err != nil {
				problems["isPrivate"] = "must be 'true' or 'false'"
			} else {
				isPrivate = &isParsed
			}
		}

		if len(problems) > 0 {
			return responses.NewValidationError("Validation error", problems)
		}

		var thumbnailData []byte
		var thumbnailFilename string

		file, fileHeader, err := r.FormFile("thumbnail")
		if err == nil {
			defer file.Close()

			if fileHeader.Size > maxFileSize {
				probs := map[string]string{
					"thumbnail": "file size must not exceed 2MB",
				}
				return responses.NewValidationError("Validation error", probs)
			}

			thumbnailData = make([]byte, fileHeader.Size)
			if _, err := file.Read(thumbnailData); err != nil {
				log.Error("failed to read thumbnail file", slog.String("error", err.Error()))
				return responses.NewApiError("File upload error", err)
			}
			thumbnailFilename = fileHeader.Filename
		} else if err != http.ErrMissingFile {
			log.Debug("unexpected file error", slog.String("error", err.Error()))
			return responses.NewApiError("File upload error", apperrors.ErrInvalidInput)
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		updatedRoom, err := updater.Update(r.Context(), id, name, categoryName, isPrivate, thumbnailData, thumbnailFilename, userId)
		if err != nil {
			return responses.NewApiError("Update room error", err)
		}

		respData := &response{
			Id:           updatedRoom.Id,
			Name:         updatedRoom.Name,
			CategoryName: updatedRoom.CategoryName,
			IsPrivate:    updatedRoom.IsPrivate,
			Thumbnail:    updatedRoom.Thumbnail,
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[response]{
			Success: true,
			Message: "Room updated",
			Data:    respData,
		})
		return nil
	}
}
