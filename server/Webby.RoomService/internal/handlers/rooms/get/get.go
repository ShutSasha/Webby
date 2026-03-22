package get

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type Getter interface {
	GetById(id uuid.UUID) (*models.Room, error)
}

func New(logger *slog.Logger, getter Getter) http.Handler {
	return errorWrapper.MakeHandler(logger, getRoom(logger, getter))
}

// getRoom godoc
// @Summary      Get a room by ID
// @Description  Retrieve a room by its ID.
// @Description  **Access Rules:**
// @Description  * **Public rooms** can be retrieved by any authenticated user
// @Description  * **Private rooms** can only be retrieved by the room creator or administrators
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Rooms
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Success      200 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Access denied - private room and not owner"
// @Failure      404 {object} docs.Error404Response "Room not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id} [get]
func getRoom(logger *slog.Logger, getter Getter) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.get"))

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
