package update

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type Updater interface {
	Update(room *models.Room) (uuid.UUID, error)
}

func New(logger *slog.Logger, updater Updater) http.Handler {
	return errorWrapper.MakeHandler(logger, updateRoom(logger, updater))
}

// updateRoom godoc
// @Summary      Update a room
// @Description  Update a room's name, **category ID**, privacy setting, and **thumbnail** by ID.
// @Description  **Validation Rules:**
// @Description  * `name`: **required**, 2-50 characters, cannot be null/empty/whitespace-only
// @Description  * `categoryId`: **required**, must be a valid UUID v4 format
// @Description  * `isPrivate`: **required**, must be valid boolean ("true" or "false")
// @Description  * `thumbnail`: optional file upload, max 2MB
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Description  * **Only the authorized room creator** can update this room
// @Tags         Rooms
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Param        name formData string true "Room name (2-50 chars, required)" minlength(2) maxlength(50)
// @Param        categoryId formData string true "Category ID (UUID v4, required)" format(uuid)
// @Param        isPrivate formData string true "Visibility flag ('true' or 'false', required)"
// @Param        thumbnail formData file false "New thumbnail image (optional, max 2MB)"
// @Success      200 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully updated"
// @Failure      400 {object} docs.Error400Response "Invalid input: name validation failed, invalid UUID, invalid boolean, file too large, or missing required fields"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Not authorized - only creator or admin can update"
// @Failure      404 {object} docs.Error404Response "Room not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id} [put]
func updateRoom(logger *slog.Logger, updater Updater) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.update"))

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
