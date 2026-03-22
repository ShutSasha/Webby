package delete

import (
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"

	"github.com/google/uuid"
)

type Deleter interface {
	Delete(id uuid.UUID) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteCategory(logger, deleter))
}

// deleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a room category by ID. **Permanent action - cannot be undone.**
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * **ONLY ADMINISTRATORS** are authorized to perform this action.
// @Tags         RoomCategory
// @Produce      json
// @Param        id path string true "Category ID"
// @Success      204 "Category successfully deleted - No Content returned"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Not authorized - Admin role required"
// @Failure      404 {object} docs.Error404Response "Category not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func deleteCategory(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.categories.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		if err := deleter.Delete(id); err != nil {
			return responses.NewApiError("Delete category error", err)
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
