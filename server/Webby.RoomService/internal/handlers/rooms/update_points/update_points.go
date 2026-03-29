package updatePoints

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

type PointsUpdater interface {
	UpdateMemberPoints(ctx context.Context, roomId uuid.UUID, memberId uuid.UUID, delta int, userId uuid.UUID) error
}

type Request struct {
	Points int `json:"points" validate:"required"`
}

func New(logger *slog.Logger, pointsUpdater PointsUpdater) http.Handler {
	return errorWrapper.MakeHandler(logger, updatePoints(logger, pointsUpdater))
}

// @Title Update member points
// @Description Add or subtract points for a room member. Only the room host can update points.
// @Param  id        path  string   true  "Room ID (UUID v4 format)"
// @Param  memberId  path  string   true  "Member User ID (UUID v4 format)"
// @Param  body      body  Request  true  "Points delta (positive to add, negative to subtract)"
// @Success  204  "Points successfully updated"
// @Failure  400  object docs.ErrorResponse  "Invalid input"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - only host can update points"
// @Failure  404  object docs.ErrorResponse  "Room or member not found"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource Rooms
// @Route /api/rooms/{id}/members/{memberId}/points [patch]
func updatePoints(logger *slog.Logger, pointsUpdater PointsUpdater) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.rooms.updatePoints"))

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

		body, problems, err := render.DecodeValid[Request](r)
		if err != nil {
			log.Debug("failed to decode request body", slog.String("error", err.Error()))
			return responses.NewApiError("Validation error", apperrors.ErrInvalidInput)
		}
		if problems != nil {
			return responses.NewValidationError("Validation error", problems)
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := pointsUpdater.UpdateMemberPoints(r.Context(), roomId, memberId, body.Points, userId); err != nil {
			return responses.NewApiError("Update points error", err)
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
