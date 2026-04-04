package create

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

type Creator interface {
	CreateVote(ctx context.Context, roomId, userId uuid.UUID, voteType, voteText string, durationSeconds int, choices []services.CreateChoiceInput) (*services.VoteDetail, error)
}

func New(logger *slog.Logger, creator Creator) http.Handler {
	return errorWrapper.MakeHandler(logger, createVote(logger, creator))
}

// @Title Create a vote
// @Description Create a new vote (poll or next_video) in a room. Only the room host can create votes.
// @Param  id    path  string                    true  "Room ID (UUID v4 format)"
// @Param  body  body  docs.CreateVoteRequest    true  "Vote creation payload"
// @Success  201  object docs.VoteApiResponse     "Vote successfully created"
// @Failure  400  object docs.ErrorResponse       "Invalid input"
// @Failure  401  object docs.ErrorResponse       "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse       "Access denied - only room host can create votes"
// @Failure  404  object docs.ErrorResponse       "Room not found"
// @Failure  500  object docs.ErrorResponse       "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes [post]
func createVote(logger *slog.Logger, creator Creator) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.create"))

	type choiceRequest struct {
		Name        string  `json:"name" validate:"required,min=1,max=200"`
		IsCorrect   bool    `json:"isCorrect"`
		QueueItemId *string `json:"queueItemId" validate:"omitempty,uuid"`
	}

	type request struct {
		Type            string          `json:"type" validate:"required,oneof=poll next_video"`
		VoteText        string          `json:"voteText" validate:"required,min=1,max=500"`
		DurationSeconds int             `json:"durationSeconds" validate:"required,min=60,max=604800"`
		Choices         []choiceRequest `json:"choices" validate:"required,min=2,max=20,dive"`
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

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		choices := make([]services.CreateChoiceInput, len(body.Choices))
		for i, c := range body.Choices {
			var queueItemId *uuid.UUID
			if c.QueueItemId != nil {
				parsed, err := uuid.Parse(*c.QueueItemId)
				if err == nil {
					queueItemId = &parsed
				}
			}
			choices[i] = services.CreateChoiceInput{
				Name:        c.Name,
				IsCorrect:   c.IsCorrect,
				QueueItemId: queueItemId,
			}
		}

		detail, err := creator.CreateVote(r.Context(), roomId, userId, body.Type, body.VoteText, body.DurationSeconds, choices)
		if err != nil {
			log.Error("Create vote error", slog.String("err", err.Error()))
			return responses.NewApiError("Create vote error", err)
		}

		render.Encode(w, r, http.StatusCreated, responses.ApiResponse[services.VoteDetail]{
			Success: true,
			Message: "Vote created",
			Data:    detail,
		})
		return nil
	}
}
