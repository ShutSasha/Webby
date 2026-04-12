package listPublic

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/models"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type PublicLister interface {
	ListPublic(ctx context.Context, page int, limit int, search string, categoryName *string) ([]models.PublicRoom, int64, error)
}

func New(logger *slog.Logger, publicLister PublicLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listPublicRooms(logger, publicLister))
}

// @Title List all public rooms
// @Description Retrieve a paginated list of all publicly visible rooms. Supports optional search and category filter.
// @Param  page      query  int     false  "Page number (default: 1)"
// @Param  limit     query  int     false  "Items per page (default: 10, max: 100)"
// @Param  search    query  string  false  "Search rooms by name (case-insensitive)"
// @Param  category  query  string  false  "Filter by category name"
// @Success  200  object docs.RoomListApiResponse  "Public rooms successfully retrieved"
// @Failure  400  object docs.ErrorResponse  "Invalid query parameters"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/public [get]
func listPublicRooms(logger *slog.Logger, publicLister PublicLister) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.listPublicRooms"))

	type roomListItem struct {
		Id            uuid.UUID `json:"id"`
		Name          string    `json:"name"`
		CategoryName  string    `json:"categoryName"`
		IsPrivate     bool      `json:"isPrivate"`
		HostId        uuid.UUID `json:"hostId"`
		HostUsername  string    `json:"hostUsername"`
		HostAvatarUrl string    `json:"hostAvatarUrl"`
		Thumbnail     string    `json:"thumbnail"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
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

		var categoryName *string
		if categoryStr := params.Get("category"); categoryStr != "" {
			if strings.ToLower(categoryStr) == "all" {
				categoryName = nil
			} else {
				categoryName = &categoryStr
			}
		}

		rooms, total, err := publicLister.ListPublic(ctx, page, limit, search, categoryName)
		if err != nil {
			log.Error("List public rooms error", slog.String("err", err.Error()))
			return err
		}

		items := make([]roomListItem, len(rooms))
		for i, room := range rooms {
			items[i] = roomListItem{
				Id:            room.Id,
				Name:          room.Name,
				CategoryName:  room.CategoryName,
				IsPrivate:     room.IsPrivate,
				HostId:        room.HostId,
				HostUsername:  room.HostUsername,
				HostAvatarUrl: room.HostAvatarUrl,
				Thumbnail:     room.Thumbnail,
			}
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[responses.PaginatedResponse[roomListItem]]{
			Success: true,
			Message: "Public rooms retrieved successfully",
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
