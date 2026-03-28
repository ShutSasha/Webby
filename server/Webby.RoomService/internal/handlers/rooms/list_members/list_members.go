package listMembers

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

type MemberLister interface {
	ListMembers(roomId uuid.UUID, page int, limit int, search string) ([]models.RoomMemberInfo, int64, error)
}

func New(logger *slog.Logger, memberLister MemberLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listMembers(logger, memberLister))
}

// listMembers godoc
// @Summary      List room members
// @Description  Retrieve a paginated list of members for a specific room.
// @Description  **Query Parameter Validation:**
// @Description  * `page`: **optional** (default: 1), must be >= 1, must be numeric
// @Description  * `limit`: **optional** (default: 10), must be 1-100, must be numeric
// @Description  * `search`: **optional**, filters members by username (case-insensitive partial match)
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Rooms
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Param        page query int false "Page number (default: 1)" minimum(1)
// @Param        limit query int false "Items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Param        search query string false "Search members by username (optional, case-insensitive)"
// @Success      200 {object} docs.ApiResponse[PaginatedResponse[RoomMemberItem]] "Room members successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format or query parameters"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id}/members [get]
func listMembers(logger *slog.Logger, memberLister MemberLister) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.listMembers"))

	type memberItem struct {
		UserId     uuid.UUID `json:"userId"`
		Username   string    `json:"username"`
		AvatarUrl  string    `json:"avatarUrl"`
		RoomPoints int       `json:"roomPoints"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Error("Invalid id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"id": "the id format is not valid",
			})
		}

		params := r.URL.Query()

		page := 1
		if pageStr := params.Get("page"); pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err != nil || p < 1 {
				return responses.NewValidationError("Validation failed", map[string]string{
					"page": "must be a positive integer",
				})
			}
			page = p
		}

		limit := 10
		if limitStr := params.Get("limit"); limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil || l < 1 || l > 100 {
				return responses.NewValidationError("Validation failed", map[string]string{
					"limit": "must be an integer between 1 and 100",
				})
			}
			limit = l
		}

		search := params.Get("search")

		members, total, err := memberLister.ListMembers(roomId, page, limit, search)
		if err != nil {
			log.Error("List room members error", slog.String("err", err.Error()))
			return responses.NewApiError("List room members error", err)
		}

		items := make([]memberItem, len(members))
		for i, member := range members {
			items[i] = memberItem{
				UserId:     member.UserId,
				Username:   member.Username,
				AvatarUrl:  member.AvatarUrl,
				RoomPoints: member.RoomPoints,
			}
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[responses.PaginatedResponse[memberItem]]{
			Success: true,
			Message: "Room members retrieved successfully",
			Data: &responses.PaginatedResponse[memberItem]{
				Items: items,
				Page:  page,
				Limit: limit,
				Total: int(total),
			},
		})

		return nil
	}
}
