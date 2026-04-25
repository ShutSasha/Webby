package get

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby-chat/internal/handlers/docs"
	errorWrapper "webby-chat/internal/handlers/errors"
	"webby-chat/internal/handlers/responses"
	"webby-chat/internal/models"
	"webby-chat/pkg/http/render"
	"webby-chat/pkg/logger"

	"github.com/google/uuid"
)

type Getter interface {
	GetById(ctx context.Context, chatId uuid.UUID) (*models.Chat, error)
}

func New(getter Getter) http.Handler {
	return errorWrapper.MakeHandler(getChat(getter))
}

// @Title Get a chat by ID
// @Description Retrieve a chat by its ID.
// @Param  id  path  string  true  "Chat ID (UUID v4 format)"
// @Success  200  object docs.ChatApiResponse  "Chat successfully retrieved"
// @Failure  400  object docs.ErrorResponse  "Invalid UUID format in path"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  404  object docs.ErrorResponse  "Chat not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Chats
// @Route /api/chats/{id} [get]
func getChat(getter Getter) errorWrapper.APIFunc {
	type chatResponse struct {
		Id        uuid.UUID  `json:"id"`
		RoomId    *uuid.UUID `json:"roomId,omitempty"`
		CreatedAt string     `json:"createdAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		log := logger.FromContext(r.Context()).With(
			slog.String("operation", "httpserver.chats.get"),
		)

		chatIdStr := r.PathValue("id")
		chatId, err := uuid.Parse(chatIdStr)
		if err != nil {
			log.Error("Invalid id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"id": "the id format is not valid",
			})
		}

		chat, err := getter.GetById(r.Context(), chatId)
		if err != nil {
			log.Error("Retrieve chat error", slog.String("err", err.Error()))
			return responses.NewApiError("Get chat error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[chatResponse]{
			Success: true,
			Message: "Chat retrieved",
			Data: &chatResponse{
				Id:        chat.Id,
				RoomId:    chat.RoomId,
				CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		})
		return nil
	}
}
