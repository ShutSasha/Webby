package create

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

type Creator interface {
	Create(ctx context.Context, room *models.Room, thumbnailData []byte, thumbnailFilename string) (*models.Room, error)
}

func New(logger *slog.Logger, creator Creator) http.Handler {
	return errorWrapper.MakeHandler(logger, createRoom(logger, creator))
}

// createRoom godoc
// @Summary      Create a new room
// @Description  Create a new room with a name, **category ID**, privacy setting, and optional **thumbnail**.
// @Description  **Validation Rules:**
// @Description  * `name`: **required**, 2-50 characters, cannot be null/empty/whitespace-only
// @Description  * `categoryId`: **required**, must be a valid UUID v4 format
// @Description  * `isPrivate`: **required**, must be valid boolean ("true" or "false")
// @Description  * `thumbnail`: optional file upload, max 2MB
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Rooms
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Room name (2-50 chars, required, non-whitespace)" minlength(2) maxlength(50)
// @Param        categoryId formData string true "Category ID (required)" format(uuid)
// @Param        isPrivate formData string true "Visibility flag (required): 'true' or 'false'"
// @Param        thumbnail formData file false "Room thumbnail image (optional, max 2MB)"
// @Success      201 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully created"
// @Failure      400 {object} docs.Error400Response "Invalid input: name empty/too short/too long, invalid categoryId UUID, invalid isPrivate boolean, file too large or missing required fields"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms [post]
func createRoom(logger *slog.Logger, creator Creator) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.create"))

	const maxFileSize = 2 * 1024 * 1024

	type roomResponse struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		CategoryId uuid.UUID `json:"categoryId"`
		IsPrivate  bool      `json:"isPrivate"`
		Thumbnail  string    `json:"thumbnail,omitempty"`
		Token      string    `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
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

		hostIdStr := r.Context().Value("userID").(string)
		hostId, _ := uuid.Parse(hostIdStr)

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

			if fileHeader.Size == 0 {
				problems := map[string]string{
					"thumbnail": "file cannot be empty",
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

		room, err := creator.Create(r.Context(), &models.Room{
			HostId:     hostId,
			CategoryId: categoryId,
			Name:       name,
			IsPrivate:  isPrivate,
		}, thumbnailData, thumbnailFilename)
		if err != nil {
			return responses.NewApiError("Create room error", err)
		}

		render.Encode(w, r, http.StatusCreated, responses.ApiResponse[roomResponse]{
			Success: true,
			Message: "Room created",
			Data: &roomResponse{
				Id:         room.Id,
				Name:       room.Name,
				CategoryId: room.CategoryId,
				IsPrivate:  room.IsPrivate,
				Thumbnail:  room.Thumbnail,
				Token:      room.Token,
			},
		})
		return nil
	}
}
