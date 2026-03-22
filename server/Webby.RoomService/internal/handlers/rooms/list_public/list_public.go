package listPublic

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"
)

type PublicLister interface {
	ListPublic(page int, limit int) ([]models.Room, int64, error)
}

func New(logger *slog.Logger, publicLister PublicLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listPublicRooms(logger, publicLister))
}

// listPublicRooms godoc
// @Summary      List all public rooms
// @Description  Retrieve a paginated list of all publicly visible rooms in the system.
// @Description  **Query Parameter Validation:**
// @Description  * `page`: **optional** (default: 1), must be >= 1, must be numeric
// @Description  * `limit`: **optional** (default: 10), must be 1-100, must be numeric
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Description  * Returns **only** public rooms (private rooms are excluded)
// @Description  * Pagination prevents excessive data retrieval
// @Tags         Rooms
// @Produce      json
// @Param        page query int false "Page number (default: 1)" minimum(1)
// @Param        limit query int false "Items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Success      200 {object} docs.ApiResponse[PaginatedResponse[RoomListItem]] "Public rooms successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid query parameters: page or limit not numeric, or outside valid range"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/public [get]
func listPublicRooms(logger *slog.Logger, publicLister PublicLister) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.listPublicRooms"))

	return func(w http.ResponseWriter, r *http.Request) error {

		return nil
	}
}
