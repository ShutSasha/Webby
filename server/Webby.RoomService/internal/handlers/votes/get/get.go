package get

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

type Getter interface {
	GetVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
}

func New(logger *slog.Logger, getter Getter) http.Handler {
	return errorWrapper.MakeHandler(logger, getVote(logger, getter))
}

// @Title Get vote details
// @Description Get detailed information about a specific vote including choices, percentages, and user's selection.
// @Param  id      path  string  true  "Room ID (UUID v4 format)"
// @Param  voteId  path  string  true  "Vote ID (UUID v4 format)"
// @Success  200  object docs.VoteApiResponse  "Vote details retrieved"
// @Failure  400  object docs.ErrorResponse    "Invalid ID format"
// @Failure  401  object docs.ErrorResponse    "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse    "Access denied - only room members"
// @Failure  404  object docs.ErrorResponse    "Vote not found"
// @Failure  500  object docs.ErrorResponse    "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes/{voteId} [get]
func getVote(logger *slog.Logger, getter Getter) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.get"))

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

		detail, err := getter.GetVote(r.Context(), voteId, userId)
		if err != nil {
			log.Error("Get vote error", slog.String("err", err.Error()))
			return responses.NewApiError("Get vote error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[services.VoteDetail]{
			Success: true,
			Message: "Vote retrieved",
			Data:    detail,
		})
		return nil
	}
}
