package httpserver

import (
	"log/slog"
	"net/http"
	"webby/internal/config"
	"webby/pkg/http/middleware"
	"webby/pkg/http/middleware/cors"
	loggerMw "webby/pkg/http/middleware/logger"
)

func NewServer(
	config *config.Config,
	logger *slog.Logger,
	roomService RoomService,
	voteService VoteService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
		logger,
		roomService,
		voteService,
	)

	var handler http.Handler = mux
	handler = middleware.Chain(handler, cors.CORS, loggerMw.New(logger))

	return handler
}
