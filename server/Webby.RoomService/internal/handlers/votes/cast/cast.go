package cast

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

type Caster interface {
	CastVote(ctx context.Context, voteId, choiceId, userId uuid.UUID) (*services.VoteDetail, error)
}

func New(logger *slog.Logger, caster Caster) http.Handler {
	return errorWrapper.MakeHandler(logger, castVote(logger, caster))
}

// @Title Cast a vote
// @Description Cast a vote for a specific choice. Only room members can vote. Each user can vote once per vote.
// @Param  id      path  string                  true  "Room ID (UUID v4 format)"
// @Param  voteId  path  string                  true  "Vote ID (UUID v4 format)"
// @Param  body    body  docs.CastVoteRequest    true  "Vote choice"
// @Success  200  object docs.VoteApiResponse    "Vote cast successfully"
// @Failure  400  object docs.ErrorResponse      "Invalid input or vote expired"
// @Failure  401  object docs.ErrorResponse      "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse      "Access denied - only room members"
// @Failure  404  object docs.ErrorResponse      "Vote or choice not found"
// @Failure  409  object docs.ErrorResponse      "User has already voted"
// @Failure  500  object docs.ErrorResponse      "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes/{voteId}/cast [post]
func castVote(logger *slog.Logger, caster Caster) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.cast"))

	type request struct {
		ChoiceId string `json:"choiceId" validate:"required,uuid"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		voteIdStr := r.PathValue("voteId")
		voteId, err := uuid.Parse(voteIdStr)
		if err != nil {
			log.Error("Invalid vote id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"voteId": "the vote id format is not valid",
			})
		}

		body, problems, err := render.DecodeValid[request](r)
		if err != nil {
			if problems != nil {
				return responses.NewValidationError("Validation failed", problems)
			}
			return responses.NewApiError("Invalid request body", err)
		}

		choiceId, _ := uuid.Parse(body.ChoiceId)

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		detail, err := caster.CastVote(r.Context(), voteId, choiceId, userId)
		if err != nil {
			log.Error("Cast vote error", slog.String("err", err.Error()))
			return responses.NewApiError("Cast vote error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[services.VoteDetail]{
			Success: true,
			Message: "Vote cast successfully",
			Data:    detail,
		})
		return nil
	}
}
