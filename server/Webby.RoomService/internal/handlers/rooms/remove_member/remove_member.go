package removeMember

import (
	"context"
	"log/slog"
	"net/http"
	"webby/internal/apperrors"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"

	"github.com/google/uuid"
)

type MemberRemover interface {
	RemoveMember(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, userId uuid.UUID) error
}

func New(logger *slog.Logger, memberRemover MemberRemover) http.Handler {
	return errorWrapper.MakeHandler(logger, removeMember(logger, memberRemover))
}

// @Title Remove a member from a room
// @Description Remove a user from a room's member list. Only the room host can remove members. The host cannot remove themselves.
// @Param  id        path  string  true  "Room ID (UUID v4 format)"
// @Param  memberId  path  string  true  "Member User ID (UUID v4 format)"
// @Success  204  "Member successfully removed"
// @Failure  400  object docs.ErrorResponse  "Invalid UUID format"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - only host can remove members"
// @Failure  404  object docs.ErrorResponse  "Room or member not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id}/members/{memberId} [delete]
func removeMember(logger *slog.Logger, memberRemover MemberRemover) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.removeMember"))

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Debug("invalid room UUID format", slog.String("id", roomIdStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		memberIdStr := r.PathValue("memberId")
		memberId, err := uuid.Parse(memberIdStr)
		if err != nil {
			log.Debug("invalid member UUID format", slog.String("memberId", memberIdStr))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := memberRemover.RemoveMember(r.Context(), roomId, memberId, userId); err != nil {
			return responses.NewApiError("Remove member error", err)
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
