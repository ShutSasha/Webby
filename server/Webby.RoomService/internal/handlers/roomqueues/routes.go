package roomqueues

import (
	"log/slog"
	"net/http"
	"webby/internal/handlers/roomqueues/add"
	"webby/internal/handlers/roomqueues/delete"
	"webby/internal/handlers/roomqueues/list"
	"webby/pkg/http/middleware/auth"
)

type QueueItemService interface {
	add.Adder
	list.QueueLister
	delete.Deleter
}

func RegisterRoomQueue(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger,
	queueItemService QueueItemService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("POST /api/rooms/{id}/queue", requireAuth(add.New(logger, queueItemService)))
	mux.Handle("GET /api/rooms/{id}/queue", requireAuth(list.New(logger, queueItemService)))
	mux.Handle("DELETE /api/rooms/{id}/queue/{itemId}", requireAuth(delete.New(logger, queueItemService)))
}
