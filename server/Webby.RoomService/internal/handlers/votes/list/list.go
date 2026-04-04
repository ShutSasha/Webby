package list

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

type Lister interface {
	ListVotes(ctx context.Context, roomId, userId uuid.UUID) ([]services.VoteDetail, error)
}

func New(logger *slog.Logger, lister Lister) http.Handler {
	return errorWrapper.MakeHandler(logger, listVotes(logger, lister))
}

// @Title List votes for a room
// @Description List all votes (polls and next_video votes) for a room. Only room members can view.
// @Param  id  path  string  true  "Room ID (UUID v4 format)"
// @Success  200  object docs.VoteListApiResponse  "Votes retrieved"
// @Failure  400  object docs.ErrorResponse         "Invalid room ID"
// @Failure  401  object docs.ErrorResponse         "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse         "Access denied - only room members"
// @Failure  500  object docs.ErrorResponse         "Internal server error"
// @Resource Votes
// @Route /api/rooms/{id}/votes [get]
func listVotes(logger *slog.Logger, lister Lister) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.votes.list"))

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Error("Invalid room id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"id": "the room id format is not valid",
			})
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		votes, err := lister.ListVotes(r.Context(), roomId, userId)
		if err != nil {
			log.Error("List votes error", slog.String("err", err.Error()))
			return responses.NewApiError("List votes error", err)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[[]services.VoteDetail]{
			Success: true,
			Message: "Votes retrieved",
			Data:    &votes,
		})
		return nil
	}
}
