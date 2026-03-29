package unvote

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/services"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Unvoter interface {
	RemoveVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
}

func New(logger *slog.Logger, unvoter Unvoter) http.Handler {
	return errorWrapper.MakeHandler(logger, removeVote(logger, unvoter))
}

// @Title Remove a user's vote
// @Description Remove the current user's vote from a vote. Only possible while the vote is still active.
// @Param  id      path  string  true  "Room ID (UUID v4 format)"
// @Param  voteId  path  string  true  "Vote ID (UUID v4 format)"
// @Success  200  object docs.VoteApiResponse  "Vote removed successfully"
// @Failure  400  object docs.ErrorResponse    "Invalid ID or vote expired"
// @Failure  401  object docs.ErrorResponse    "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse    "Access denied - only room members"
// @Failure  404  object docs.ErrorResponse    "Vote not found or user hasn't voted"
// @Failure  500  object docs.ErrorResponse    "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes/{voteId}/cast [delete]
func removeVote(logger *slog.Logger, unvoter Unvoter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.unvote"))

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

		detail, err := unvoter.RemoveVote(r.Context(), voteId, userId)
		if err != nil {
			log.Error("Remove vote error", slog.String("err", err.Error()))
			return responses.NewApiError("Remove vote error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[services.VoteDetail]{
			Success: true,
			Message: "Vote removed successfully",
			Data:    detail,
		})
		return nil
	}
}
