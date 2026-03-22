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

// updateCategory godoc
// @Summary      Update a category
// @Description  Update an existing room category by ID. Only the **name** field is updateable.
// @Description  **Validation Rules:**
// @Description  * Name is **required** (cannot be null, empty, or whitespace-only)
// @Description  * Name length must be **2-50 characters**
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * **ONLY ADMINISTRATORS** are authorized to perform this action.
// @Tags         RoomCategory
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID (UUID v4 format)"
// @Param        request body docs.UpdateCategoryRequest true "Category update data - Name (2-50 chars, required, non-whitespace)"
// @Success      200 {object} docs.ApiResponse[docs.CategoryResponse] "Category successfully updated"
// @Failure      400 {object} docs.Error400Response "Invalid input: name missing, empty, too short (<2), too long (>50), whitespace-only, invalid UUID in path, or malformed JSON"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Not authorized - Admin role required"
// @Failure      404 {object} docs.Error404Response "Category not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /categories/{id} [put]
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
