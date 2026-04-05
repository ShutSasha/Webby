package addMember

import (
	"context"
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type MemberAdder interface {
	AddMembers(ctx context.Context, roomId uuid.UUID, memberIds []uuid.UUID, userId uuid.UUID) error
}

func New(logger *slog.Logger, memberAdder MemberAdder) http.Handler {
	return errorWrapper.MakeHandler(logger, addMember(logger, memberAdder))
}

// @Title Add members to a room
// @Description Add users to a room's member list. Only the room host can add members.
// @Param  id    path  string                  true  "Room ID (UUID v4 format)"
// @Param  body  body  docs.AddMembersRequest  true  "User IDs to add"
// @Success  204  "Members successfully added"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - only host can add members"
// @Failure  404  object docs.ErrorResponse  "Room not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id}/members [post]
func addMember(logger *slog.Logger, memberAdder MemberAdder) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.addMember"))

	type request struct {
		UserIds []string `json:"userIds" validate:"required,min=1,dive,uuid"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Debug("invalid room UUID format", slog.String("id", roomIdStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		body, problems, err := render.DecodeValid[request](r)
		if err != nil {
			log.Debug("failed to decode request body", slog.String("error", err.Error()))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}
		if problems != nil {
			return responses.NewValidationError("Validation error", problems)
		}

		memberIds := make([]uuid.UUID, len(body.UserIds))
		for i, id := range body.UserIds {
			memberIds[i], _ = uuid.Parse(id)
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := memberAdder.AddMembers(r.Context(), roomId, memberIds, userId); err != nil {
			return responses.NewApiError("Add member error", err)
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
