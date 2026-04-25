package httpserver

import (
	"log/slog"
	"net/http"
	"webby-chat/internal/config"
	"webby-chat/pkg/http/middleware"
	"webby-chat/pkg/http/middleware/cors"
	loggerMw "webby-chat/pkg/http/middleware/logger"

	socketio "github.com/googollee/go-socket.io"
)

func NewServer(
	config *config.Config,
	logger *slog.Logger,
	chatService ChatService,
	socketServer *socketio.Server,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
		chatService,
		socketServer,
	)

	var handler http.Handler = mux
	handler = middleware.Chain(handler, cors.CORS, loggerMw.New(logger))

	return handler
}
