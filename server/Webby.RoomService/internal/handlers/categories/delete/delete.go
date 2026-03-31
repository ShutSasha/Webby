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

// @Title Delete a category
// @Description Delete a room category by ID. Permanent action. Only administrators can perform this action.
// @Param  id  path  string  true  "Category ID (UUID v4 format)"
// @Success  204  "Category successfully deleted"
// @Failure  400  object docs.ErrorResponse  "Invalid UUID format in path"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - Admin role required"
// @Failure  404  object docs.ErrorResponse  "Category not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomCategory
// @Route /api/categories/{id} [delete]
func deleteCategory(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.categories.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Debug("Validation error", slog.Any("err", err.Error()))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		if err := deleter.Delete(id); err != nil {
			log.Debug("Delete category error", slog.Any("err", err.Error()))
			return responses.NewApiError("Delete category error", err)
		}

		log.Debug("Category successfully deleted", slog.String("id", idStr))
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
