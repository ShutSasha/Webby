package httpserver

import (
	"net/http"
	"webby-chat/internal/config"
	"webby-chat/internal/handlers/chats"
	"webby-chat/internal/handlers/chats/create"
	"webby-chat/internal/handlers/chats/get"

	socketio "github.com/googollee/go-socket.io"
)

type ChatService interface {
	create.Creator
	get.Getter
}

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	chatService ChatService,
	socketServer *socketio.Server,
) {
	mux.Handle("/api/", http.NotFoundHandler())

	chats.RegisterChats(mux, []byte(cfg.JwtSecret), chatService)

	mux.Handle("/socket.io/", socketServer)
}
