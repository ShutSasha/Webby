package votes

import (
	"log/slog"
	"net/http"
	"webby/internal/handlers/votes/cast"
	"webby/internal/handlers/votes/create"
	voteDelete "webby/internal/handlers/votes/delete"
	"webby/internal/handlers/votes/get"
	"webby/internal/handlers/votes/list"
	"webby/internal/handlers/votes/unvote"
	"webby/pkg/http/middleware/auth"
)

type VoteService interface {
	create.Creator
	list.Lister
	get.Getter
	cast.Caster
	unvote.Unvoter
	voteDelete.Deleter
}

func RegisterVotes(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger,
	voteService VoteService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("POST /api/rooms/{id}/votes", requireAuth(create.New(logger, voteService)))
	mux.Handle("GET /api/rooms/{id}/votes", requireAuth(list.New(logger, voteService)))
	mux.Handle("GET /api/rooms/{id}/votes/{voteId}", requireAuth(get.New(logger, voteService)))
	mux.Handle("POST /api/rooms/{id}/votes/{voteId}/cast", requireAuth(cast.New(logger, voteService)))
	mux.Handle("DELETE /api/rooms/{id}/votes/{voteId}/cast", requireAuth(unvote.New(logger, voteService)))
	mux.Handle("DELETE /api/rooms/{id}/votes/{voteId}", requireAuth(voteDelete.New(logger, voteService)))
}
