package chats

import (
	"log/slog"
	"net/http"
	"webby-chat/internal/handlers/chats/create"
	"webby-chat/internal/handlers/chats/get"
	"webby-chat/pkg/http/middleware/auth"
)

type ChatService interface {
	create.Creator
	get.Getter
}

func RegisterChats(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger, chatService ChatService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("POST /api/chats", requireAuth(create.New(logger, chatService)))
	mux.Handle("GET /api/chats/{id}", requireAuth(get.New(logger, chatService)))
}
