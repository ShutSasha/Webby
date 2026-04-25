package chats

import (
	"net/http"
	"webby-chat/internal/handlers/chats/create"
	"webby-chat/internal/handlers/chats/get"
	"webby-chat/pkg/http/middleware/auth"
)

type ChatService interface {
	create.Creator
	get.Getter
}

func RegisterChats(mux *http.ServeMux, jwtSecret []byte, chatService ChatService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("POST /api/chats", requireAuth(create.New(chatService)))
	mux.Handle("GET /api/chats/{id}", requireAuth(get.New(chatService)))
}
