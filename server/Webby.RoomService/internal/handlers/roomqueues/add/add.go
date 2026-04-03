package add

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

type Adder interface {
	AddToQueue(ctx context.Context, roomId, userId, entityId uuid.UUID, entityType string) (*models.QueueItem, error)
}

func New(logger *slog.Logger, adder Adder) http.Handler {
	return errorWrapper.MakeHandler(logger, addToQueue(logger, adder))
}

// @Title Add item to room queue
// @Description Add a video or playlist to the room's playback queue. Only room members can add items.
// @Param  id    path  string          true  "Room ID (UUID v4 format)"
// @Param  body  body  docs.AddQueueItemRequest     true  "Queue item to add"
// @Success  201  object docs.QueueItemApiResponse  "Item successfully added to queue"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Access denied - only room members can add to queue"
// @Failure  404  object docs.ErrorResponse  "Room or entity not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomQueue
// @Route /api/rooms/{id}/queue [post]
func addToQueue(logger *slog.Logger, adder Adder) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.roomqueues.add"))

	type request struct {
		EntityId   string `json:"entityId" validate:"required,uuid"`
		EntityType string `json:"entityType" validate:"required,oneof=video playlist"`
	}

	type queueItemResponse struct {
		Id         uuid.UUID `json:"id"`
		EntityId   uuid.UUID `json:"entityId"`
		EntityType string    `json:"entityType"`
		IsActive   bool      `json:"isActive"`
		Position   int       `json:"position"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Error("Invalid room id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"id": "the room id format is not valid",
			})
		}

		body, problems, err := render.DecodeValid[request](r)
		if err != nil {
			if problems != nil {
				return responses.NewValidationError("Validation failed", problems)
			}
			return responses.NewApiError("Invalid request body", err)
		}

		entityId, _ := uuid.Parse(body.EntityId)

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		item, err := adder.AddToQueue(r.Context(), roomId, userId, entityId, body.EntityType)
		if err != nil {
			log.Error("Add to queue error", slog.String("err", err.Error()))
			return responses.NewApiError("Add to queue error", err)
		}

		render.Encode(w, r, http.StatusCreated, responses.ApiResponse[queueItemResponse]{
			Success: true,
			Message: "Item added to queue",
			Data: &queueItemResponse{
				Id:         item.Id,
				EntityId:   item.EntityId,
				EntityType: item.EntityType,
				IsActive:   item.IsActive,
				Position:   item.Position,
			},
		})
		return nil
	}
}
