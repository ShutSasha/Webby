package list

import (
	"log/slog"
	"net/http"
	"strconv"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/models"
	"webby/pkg/http/render"

	"github.com/google/uuid"
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
// @Description  **Query Parameters Fallback:**
// @Description  * Invalid, missing, or non-numeric `page` values will default to **1**.
// @Description  * Invalid, missing, or non-numeric `limit` values (<1) will default to **10**. Limits >100 are capped at **100**.
// @Description  * `search`: optional string to filter categories by name (case-insensitive)
// @Tags         RoomCategory
// @Produce      json
// @Param        search query string false "Search categories by name (optional)"
// @Param        page query int false "Page number (defaults to 1 if invalid)" default(1)
// @Param        limit query int false "Items per page (defaults to 10 if invalid, max 100)" default(10)
// @Success      200 {object} docs.ApiResponse[docs.PaginatedResponse[docs.CategoryResponse]] "Successful retrieval with pagination metadata"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Router       /categories [get]
func listCategories(logger *slog.Logger, lister Lister) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.categories.list"))

	type categoryResponse struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		params := r.URL.Query()

		page, _ := strconv.Atoi(params.Get("page"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(params.Get("limit"))
		if limit < 1 {
			limit = 10
		} else if limit > 100 {
			limit = 100
		}

		search := params.Get("search")

		categories, total, err := lister.List(search, page, limit)
		if err != nil {
			return err
		}

		items := make([]categoryResponse, len(categories))
		for i, c := range categories {
			items[i] = categoryResponse{
				Id:   c.Id,
				Name: c.Name,
			}
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[responses.PaginatedResponse[categoryResponse]]{
			Success: true,
			Message: "Categories retrieved successfully",
			Data: &responses.PaginatedResponse[categoryResponse]{
				Items: items,
				Page:  page,
				Limit: limit,
				Total: int(total),
			},
		})

		return nil
	}
}
