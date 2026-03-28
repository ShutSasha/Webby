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

// getRoom godoc
// @Summary      Get a room by ID
// @Description  Retrieve a room by its ID.
// @Description  **Access Rules:**
// @Description  * Only the **room creator (host)** can retrieve the room by ID
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Rooms
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Success      200 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Access denied - only the room creator can retrieve"
// @Failure      404 {object} docs.Error404Response "Room not found"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id} [get]
func getRoom(logger *slog.Logger, getter Getter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.get"))

	type response struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		CategoryId uuid.UUID `json:"categoryId"`
		IsPrivate  bool      `json:"isPrivate"`
		HostId     uuid.UUID `json:"hostId"`
		Thumbnail  string    `json:"thumbnail"`
		Token      string    `json:"token"`
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
				Id:         room.Id,
				Name:       room.Name,
				CategoryId: room.CategoryId,
				IsPrivate:  room.IsPrivate,
				HostId:     room.HostId,
				Thumbnail:  room.Thumbnail,
				Token:      room.Token,
			},
		})
		return nil
	}
}
