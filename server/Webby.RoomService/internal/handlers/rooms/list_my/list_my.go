package listMy

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

type MyLister interface {
	ListMy(id uuid.UUID, page int, limit int) ([]models.Room, int64, error)
}

func New(logger *slog.Logger, myLister MyLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listMyRooms(logger, myLister))
}

// @Title List authenticated user's rooms
// @Description Retrieve a paginated list of all rooms created by the authenticated user. Returns both public and private rooms.
// @Param  page   query  int  false  "Page number (default: 1)"
// @Param  limit  query  int  false  "Items per page (default: 10, max: 100)"
// @Success  200  object docs.RoomListApiResponse  "User's rooms successfully retrieved"
// @Failure  400  object docs.ErrorResponse  "Invalid query parameters"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/my [get]
func listMyRooms(logger *slog.Logger, myLister MyLister) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.listMyRooms"))

	type roomListItem struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		CategoryId uuid.UUID `json:"categoryId"`
		IsPrivate  bool      `json:"isPrivate"`
		HostId     uuid.UUID `json:"hostId"`
		Thumbnail  string    `json:"thumbnail"`
		Token      string    `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
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

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		rooms, total, err := myLister.ListMy(userId, page, limit)
		if err != nil {
			log.Error("List my rooms error", slog.String("err", err.Error()))
			return err
		}

		items := make([]roomListItem, len(rooms))
		for i, room := range rooms {
			items[i] = roomListItem{
				Id:         room.Id,
				Name:       room.Name,
				CategoryId: room.CategoryId,
				IsPrivate:  room.IsPrivate,
				HostId:     room.HostId,
				Thumbnail:  room.Thumbnail,
				Token:      room.Token,
			}
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[responses.PaginatedResponse[roomListItem]]{
			Success: true,
			Message: "User rooms retrieved successfully",
			Data: &responses.PaginatedResponse[roomListItem]{
				Items: items,
				Page:  page,
				Limit: limit,
				Total: int(total),
			},
		})

		return nil
	}
}
