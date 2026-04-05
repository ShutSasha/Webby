package delete

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Deleter interface {
	DeleteVote(ctx context.Context, voteId, userId uuid.UUID) error
}

func New(logger *slog.Logger, deleter Deleter) http.Handler {
	return errorWrapper.MakeHandler(logger, deleteVote(logger, deleter))
}

// @Title Delete a vote
// @Description Delete a vote. Only the room host can delete votes.
// @Param  id      path  string  true  "Room ID (UUID v4 format)"
// @Param  voteId  path  string  true  "Vote ID (UUID v4 format)"
// @Success  200  object docs.SuccessResponse  "Vote deleted"
// @Failure  400  object docs.ErrorResponse    "Invalid ID format"
// @Failure  401  object docs.ErrorResponse    "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse    "Access denied - only room host"
// @Failure  404  object docs.ErrorResponse    "Vote not found"
// @Failure  500  object docs.ErrorResponse    "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes/{voteId} [delete]
func deleteVote(logger *slog.Logger, deleter Deleter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.delete"))

	return func(w http.ResponseWriter, r *http.Request) error {
		voteIdStr := r.PathValue("voteId")
		voteId, err := uuid.Parse(voteIdStr)
		if err != nil {
			log.Error("Invalid vote id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"voteId": "the vote id format is not valid",
			})
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		if err := deleter.DeleteVote(r.Context(), voteId, userId); err != nil {
			log.Error("Delete vote error", slog.String("err", err.Error()))
			return responses.NewApiError("Delete vote error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[struct{}]{
			Success: true,
			Message: "Vote deleted",
		})
		return nil
	}
}
