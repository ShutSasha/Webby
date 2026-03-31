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

// @Title List categories
// @Description Retrieve a paginated list of room categories. Supports optional search by name.
// @Param  search  query  string  false  "Search categories by name"
// @Param  page    query  int     false  "Page number (default: 1)"
// @Param  limit   query  int     false  "Items per page (default: 10, max: 100)"
// @Success  200  object docs.CategoryListApiResponse  "Categories successfully retrieved"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomCategory
// @Route /api/categories [get]
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
