package delete

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Deleter interface {
	DeleteFromQueue(ctx context.Context, itemId, userId uuid.UUID) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteFromQueue(logger, deleter))
}

// deleteFromQueue godoc
// @Summary      Remove item from room queue
// @Description  Remove a video or playlist from the room's playback queue.
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format (room ID)
// @Description  * `itemId`: must be a valid UUID v4 format (queue item ID)
// @Description  **Access Rules:**
// @Description  * Only **room members** can remove items from the queue
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Description  **PAY ATTENTION:**
// @Description  * ONLY A SEPARATE VIDEO OR THE PLAYLIST ITSELF CAN BE DELETED, THE VIDEO FROM THE PLAYLIST CAN NOT BE DELETED SEPARATELY**
// @Tags         Room Queue
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Param        itemId path string true "Queue Item ID (UUID v4 format)"
// @Success      200 {object} docs.ApiResponse[any] "Item successfully removed from queue"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Access denied - only room members can remove from queue"
// @Failure      404 {object} docs.Error404Response "Queue item not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id}/queue/{itemId} [delete]
func deleteFromQueue(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.roomqueues.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		itemIdStr := r.PathValue("itemId")
		itemId, err := uuid.Parse(itemIdStr)
		if err != nil {
			log.Error("Invalid item id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"itemId": "the item id format is not valid",
			})
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := deleter.DeleteFromQueue(r.Context(), itemId, userId); err != nil {
			log.Error("Delete from queue error", slog.String("err", err.Error()))
			return responses.NewApiError("Delete from queue error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[struct{}]{
			Success: true,
			Message: "Item removed from queue",
		})
		return nil
	}
}
