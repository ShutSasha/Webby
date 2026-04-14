package delete

import (
	"context"
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Deleter interface {
	Delete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteRoom(logger, deleter))
}

// @Title Delete a room
// @Description Delete a room by ID. Permanent action - cannot be undone. Only the authorized room creator can delete.
// @Param  id  path  string  true  "Room ID (UUID v4 format)"
// @Success  204  "Room successfully deleted"
// @Failure  400  object docs.ErrorResponse  "Invalid UUID format in path"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized"
// @Failure  404  object docs.ErrorResponse  "Room not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id} [delete]
func deleteRoom(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Debug("invalid UUID format", slog.String("id", idStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := deleter.Delete(r.Context(), id, userId); err != nil {
			return responses.NewApiError("Delete room error", err)
		}

		render.Encode(w, r, http.StatusOK, &responses.ApiResponse[any]{
			Success: true,
			Message: "Room successfully deleted",
		})
		return nil
	}
}
