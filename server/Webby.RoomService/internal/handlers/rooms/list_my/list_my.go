package listMy

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type MyLister interface {
	ListMy(id uuid.UUID, page int, limit int) ([]models.Room, int64, error)
}

func New(logger *slog.Logger, myLister MyLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listMyRooms(logger, myLister))
}

// listMyRooms godoc
// @Summary      List authenticated user's rooms
// @Description  Retrieve a paginated list of all rooms created by the authenticated user.
// @Description  **Query Parameter Validation:**
// @Description  * `page`: **optional** (default: 1), must be >= 1, must be numeric
// @Description  * `limit`: **optional** (default: 10), must be 1-100, must be numeric
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Description  * Returns **only** the authenticated user's rooms
// @Description  * List includes both public and private rooms owned by the user
// @Tags         Rooms
// @Produce      json
// @Param        page query int false "Page number (default: 1)" minimum(1)
// @Param        limit query int false "Items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Success      200 {object} docs.ApiResponse[PaginatedResponse[RoomListItem]] "User's rooms successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid query parameters: page or limit not numeric, or outside valid range"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/my [get]
func listMyRooms(logger *slog.Logger, myLister MyLister) errorWrapper.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.rooms.listMyRooms"))

	return func(w http.ResponseWriter, r *http.Request) error {

		return nil
	}
}
