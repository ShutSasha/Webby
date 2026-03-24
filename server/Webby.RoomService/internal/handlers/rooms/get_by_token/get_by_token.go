package getByToken

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/models"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type TokenGetter interface {
	GetByToken(token string) (*models.Room, error)
}

func New(logger *slog.Logger, tokenGetter TokenGetter) http.Handler {
	return errorWrapper.MakeHandler(logger, getByToken(logger, tokenGetter))
}

// getByToken godoc
// @Summary      Get a room by invite token
// @Description  Retrieve room information using its unique invite token.
// @Description  **Request Body Validation:**
// @Description  * `token`: **required**, cannot be empty or whitespace-only
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Rooms
// @Accept       json
// @Produce      json
// @Param        request body object{token=string} true "Room invite token"
// @Success      200 {object} docs.ApiResponse[docs.RoomResponse] "Room successfully retrieved by token"
// @Failure      400 {object} docs.Error400Response "Invalid request body or missing token"
// @Failure      404 {object} docs.Error404Response "Room with given token not found"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/token [post]
func getByToken(logger *slog.Logger, tokenGetter TokenGetter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.getByToken"))

	type request struct {
		Token string `json:"token" validate:"required,notblank"`
	}

	type roomResponse struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		CategoryId uuid.UUID `json:"categoryId"`
		IsPrivate  bool      `json:"isPrivate"`
		HostId     uuid.UUID `json:"hostId"`
		Thumbnail  string    `json:"thumbnail"`
		Token      string    `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		req, problems, err := render.DecodeValid[request](r)
		if len(problems) > 0 {
			log.Debug("validation problems", slog.Any("problems", problems))
			return responses.NewValidationError("Validation failed", problems)
		}
		if err != nil {
			return responses.NewApiError("Validation error", err)
		}

		room, err := tokenGetter.GetByToken(req.Token)
		if err != nil {
			log.Error("Get room by token error", slog.String("err", err.Error()))
			return responses.NewApiError("Get room by token error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[roomResponse]{
			Success: true,
			Message: "Room retrieved successfully",
			Data: &roomResponse{
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
