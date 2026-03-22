package create

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type Creator interface {
	Create(room *models.Room) (uuid.UUID, error)
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
// @Param        categoryId formData string true "Category ID (UUID v4, required)" format(uuid)
// @Param        isPrivate formData string true "Visibility flag (required): 'true' or 'false'"
// @Param        thumbnail formData file false "Room thumbnail image (optional, max 2MB)"
// @Success      201 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully created"
// @Failure      400 {object} docs.Error400Response "Invalid input: name empty/too short/too long, invalid categoryId UUID, invalid isPrivate boolean, file too large, wrong Content-Type, or missing required fields"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms [post]
func createRoom(logger *slog.Logger, creator Creator) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.create"))

	type createRoomRequest struct{}

	type roomResponse struct{}

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusCreated)
		return nil
	}
}
