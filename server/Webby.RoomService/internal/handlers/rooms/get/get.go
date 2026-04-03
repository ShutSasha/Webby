package get

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/models"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Getter interface {
	GetById(ctx context.Context, roomId, userId uuid.UUID) (*models.Room, error)
}

func New(logger *slog.Logger, getter Getter) http.Handler {
	return errorWrapper.MakeHandler(logger, getRoom(logger, getter))
}

// @Title Get a room by ID
// @Description Retrieve a room by its ID. For public rooms, the requesting user is automatically added as a member. For private rooms, only existing members can access the room.
// @Param  id  path  string  true  "Room ID (UUID v4 format)"
// @Success  200  object docs.RoomApiResponse  "Room successfully retrieved"
// @Failure  400  object docs.ErrorResponse  "Invalid UUID format in path"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - not a member of private room"
// @Failure  404  object docs.ErrorResponse  "Room not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id} [get]
func getRoom(logger *slog.Logger, getter Getter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.get"))

	type response struct {
		Id           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		CategoryName string    `json:"categoryName"`
		IsPrivate    bool      `json:"isPrivate"`
		HostId       uuid.UUID `json:"hostId"`
		Thumbnail    string    `json:"thumbnail"`
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

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		room, err := getter.GetById(r.Context(), roomId, userId)
		if err != nil {
			log.Error("Retrieve room error", slog.String("err", err.Error()))
			return responses.NewApiError("Get room error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[response]{
			Success: true,
			Message: "Room retrieved",
			Data: &response{
				Id:           room.Id,
				Name:         room.Name,
				CategoryName: room.CategoryName,
				IsPrivate:    room.IsPrivate,
				HostId:       room.HostId,
				Thumbnail:    room.Thumbnail,
			},
		})
		return nil
	}
}
