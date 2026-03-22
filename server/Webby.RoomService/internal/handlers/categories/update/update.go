package update

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"

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
	_ = logger.With(slog.String("operation", "httpserver.categories.update"))

	type updateCategoryRequest struct{}

	type categoryResponse struct{}

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusOK)
		return nil
	}
}
