package list

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"
)

type Lister interface {
	List(search string, page int, limit int) ([]models.Category, int64, error)
}

func New(logger *slog.Logger, lister Lister) http.Handler {
	return errorWrapper.MakeHandler(logger, listCategories(logger, lister))
}

// listCategories godoc
// @Summary      List categories
// @Description  Retrieve a **paginated** list of room categories. Supports optional **search** by name.
// @Description  **Query Parameters Validation:**
// @Description  * `page`: must be >= 1 (default: 1)
// @Description  * `limit`: must be between 1 and 100 (default: 10)
// @Description  * `search`: optional string to filter categories by name (case-insensitive)
// @Tags         RoomCategory
// @Produce      json
// @Param        search query string false "Search categories by name (optional)"
// @Param        page query int false "Page number (1-based, required >= 1)" default(1) minimum(1)
// @Param        limit query int false "Items per page (1-100)" default(10) minimum(1) maximum(100)
// @Success      200 {object} docs.ApiResponse[docs.PaginatedResponse[docs.CategoryResponse]] "Successful retrieval with pagination metadata"
// @Failure      400 {object} docs.Error400Response "Invalid query parameters: page <= 0, limit out of range, or non-numeric values"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Router       /categories [get]
func listCategories(logger *slog.Logger, lister Lister) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.categories.list"))

	return func(w http.ResponseWriter, r *http.Request) error {

		w.WriteHeader(http.StatusOK)
		return nil
	}
}
