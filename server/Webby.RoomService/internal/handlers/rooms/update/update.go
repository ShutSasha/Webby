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
	Update(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string, userId uuid.UUID) (uuid.UUID, error)
}

func New(logger *slog.Logger, updater Updater) http.Handler {
	return errorWrapper.MakeHandler(logger, updateRoom(logger, updater))
}

// @Title Update a room
// @Description Update a room's name, category ID, privacy setting, and thumbnail by ID. Only the authorized room creator can update.
// @Param  id          path  string  true   "Room ID (UUID v4 format)"
// @Param  name        form  string  true   "Room name (2-50 chars, required)"
// @Param  categoryId  form  string  true   "Category ID (UUID v4, required)"
// @Param  isPrivate   form  string  true   "Visibility flag ('true' or 'false', required)"
// @Param  thumbnail   file  file    false  "New thumbnail image (optional, max 2MB)"
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
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		CategoryId uuid.UUID `json:"categoryId"`
		IsPrivate  bool      `json:"isPrivate"`
		Thumbnail  string    `json:"thumbnail,omitempty"`
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

		name := r.FormValue("name")
		categoryIdStr := r.FormValue("categoryId")
		isPrivateStr := r.FormValue("isPrivate")

		if strings.TrimSpace(name) == "" || len(name) < 2 || len(name) > 50 {
			problems := map[string]string{
				"name": "must be between 2 and 50 characters and cannot be empty",
			}
			return responses.NewValidationError("Validation error", problems)
		}

		categoryId, err := uuid.Parse(categoryIdStr)
		if err != nil {
			problems := map[string]string{
				"categoryId": "invalid UUID format",
			}
			log.Debug("invalid categoryId format", slog.String("categoryId", categoryIdStr))
			return responses.NewValidationError("Validation error", problems)
		}

		isPrivate, err := strconv.ParseBool(isPrivateStr)
		if err != nil {
			problems := map[string]string{
				"isPrivate": "must be 'true' or 'false'",
			}
			log.Debug("invalid isPrivate format", slog.String("isPrivate", isPrivateStr))
			return responses.NewValidationError("Validation error", problems)
		}

		var thumbnailData []byte
		var thumbnailFilename string

		file, fileHeader, err := r.FormFile("thumbnail")
		if err == nil {
			defer file.Close()

			if fileHeader.Size > maxFileSize {
				problems := map[string]string{
					"thumbnail": "file size must not exceed 2MB",
				}
				return responses.NewValidationError("Validation error", problems)
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

		room := &models.Room{
			Id:         id,
			CategoryId: categoryId,
			Name:       name,
			IsPrivate:  isPrivate,
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		_, err = updater.Update(r.Context(), room, thumbnailData, thumbnailFilename, userId)
		if err != nil {
			return responses.NewApiError("Update room error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[response]{
			Success: true,
			Message: "Room updated",
			Data: &response{
				Id:         id,
				Name:       name,
				CategoryId: categoryId,
				IsPrivate:  isPrivate,
				Thumbnail:  room.Thumbnail,
			},
		})
		return nil
	}
}
