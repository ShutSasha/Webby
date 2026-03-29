package update

import (
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Updater interface {
	Update(id uuid.UUID, name string) error
}

func New(logger *slog.Logger, updater Updater) http.Handler {
	return errorWrapper.MakeHandler(logger, updateCategory(logger, updater))
}

// @Title Update a category
// @Description Update an existing room category by ID. Only the name field is updateable. Only administrators can perform this action.
// @Param  id       path  string          true  "Category ID (UUID v4 format)"
// @Param  request  body  docs.UpdateCategoryRequest  true  "Category update data - Name (2-50 chars, required)"
// @Success  200  object docs.CategoryApiResponse  "Category successfully updated"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - Admin role required"
// @Failure  404  object docs.ErrorResponse  "Category not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomCategory
// @Route /api/categories/{id} [put]
func updateCategory(logger *slog.Logger, updater Updater) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.categories.update"))

	type request struct {
		Name string `json:"name" validate:"required,notblank,min=2,max=50"`
	}

	type response struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Debug("invalid UUID format", slog.String("id", idStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		var req request
		req, problems, err := render.DecodeValid[request](r)
		if len(problems) > 0 {
			log.Debug("problems debug", slog.Any("problems", problems))
			return responses.NewValidationError("Validation error", problems)
		}
		if err != nil {
			return responses.NewApiError("Validation error", err)
		}

		if err := updater.Update(id, req.Name); err != nil {
			return responses.NewApiError("Update category error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[response]{
			Success: true,
			Message: "Category updated",
			Data: &response{
				Id:   id,
				Name: req.Name,
			},
		})
		return nil
	}
}
