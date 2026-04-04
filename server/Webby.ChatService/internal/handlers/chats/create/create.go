package create

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby-chat/internal/handlers/docs"
	errorWrapper "webby-chat/internal/handlers/errors"
	"webby-chat/internal/handlers/responses"
	"webby-chat/internal/models"
	"webby-chat/pkg/http/render"

	"github.com/google/uuid"
)

type Creator interface {
	Create(ctx context.Context, roomId *uuid.UUID) (*models.Chat, error)
}

func New(logger *slog.Logger, creator Creator) http.Handler {
	return errorWrapper.MakeHandler(logger, createChat(logger, creator))
}

// @Title Create a new chat
// @Description Create a new chat, optionally linked to a room.
// @Param  body  body  docs.CreateRequest  true  "Chat creation request"
// @Success  201  object docs.ChatApiResponse  "Chat successfully created"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Chats
// @Route /api/chats [post]
func createChat(logger *slog.Logger, creator Creator) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.chats.create"))

	type createRequest struct {
		RoomId *string `json:"roomId"`
	}

	type chatResponse struct {
		Id        uuid.UUID  `json:"id"`
		RoomId    *uuid.UUID `json:"roomId,omitempty"`
		CreatedAt string     `json:"createdAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := render.Decode[createRequest](r)
		if err != nil {
			log.Debug("failed to decode request", slog.String("error", err.Error()))
			return responses.NewValidationError("Validation error", map[string]string{
				"body": "invalid JSON body",
			})
		}

		var roomId *uuid.UUID
		if req.RoomId != nil && *req.RoomId != "" {
			parsed, err := uuid.Parse(*req.RoomId)
			if err != nil {
				return responses.NewValidationError("Validation error", map[string]string{
					"roomId": "invalid UUID format",
				})
			}
			roomId = &parsed
		}

		chat, err := creator.Create(r.Context(), roomId)
		if err != nil {
			log.Error("Create chat error", slog.String("err", err.Error()))
			return responses.NewApiError("Create chat error", err)
		}

		render.Encode(w, r, http.StatusCreated, responses.ApiResponse[chatResponse]{
			Success: true,
			Message: "Chat created",
			Data: &chatResponse{
				Id:        chat.Id,
				RoomId:    chat.RoomId,
				CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		})
		return nil
	}
}
