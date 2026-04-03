package delete

import (
	"context"
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
)

type Deleter interface {
	Delete(ctx context.Context, name string) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteCategory(logger, deleter))
}

// @Title Delete a category
// @Description Delete a room category by name. Permanent action. Only administrators can perform this action.
// @Param  name  path  string  true  "Category name"
// @Success  204  "Category successfully deleted"
// @Failure  400  object docs.ErrorResponse  "Invalid category name"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - Admin role required"
// @Failure  404  object docs.ErrorResponse  "Category not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomCategory
// @Route /api/categories/{name} [delete]
func deleteCategory(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.categories.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		name := r.PathValue("name")

		if name == "" {
			log.Debug("Validation error: empty category name")
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		if err := deleter.Delete(ctx, name); err != nil {
			log.Debug("Delete category error", slog.Any("err", err.Error()))
			return responses.NewApiError("Delete category error", err)
		}

		log.Debug("Category successfully deleted", slog.String("name", name))
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
