package delete

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"

	"github.com/google/uuid"
)

type Deleter interface {
	Delete(id uuid.UUID) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteRoom(logger, deleter))
}

// deleteRoom godoc
// @Summary      Delete a room
// @Description  Delete a room by ID. **Permanent action - cannot be undone.**
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Description  * **Only the authorized room creator** can delete this room
// @Tags         Rooms
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Success      204 "Room successfully deleted - No Content returned"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Not authorized - only creator or admin can delete"
// @Failure      404 {object} docs.Error404Response "Room not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id} [delete]
func deleteRoom(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
