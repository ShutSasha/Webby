package rooms

import (
	"log/slog"
	"net/http"
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/get"
	listMy "webby/internal/handlers/rooms/list_my"
	listPublic "webby/internal/handlers/rooms/list_public"
	"webby/internal/handlers/rooms/update"
	"webby/pkg/http/middleware/auth"
)

type RoomService interface {
	create.Creator
	listMy.MyLister
	listPublic.PublicLister
	get.Getter
	update.Updater
	delete.Deleter
}

func RegisterRooms(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger,
	roomService RoomService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("GET /rooms/public", listPublic.New(logger, roomService))

	mux.Handle("POST /rooms", requireAuth(create.New(logger, roomService)))
	mux.Handle("GET /rooms/my", requireAuth(listMy.New(logger, roomService)))
	mux.Handle("GET /rooms/{id}", requireAuth(get.New(logger, roomService)))
	mux.Handle("PUT /rooms/{id}", requireAuth(update.New(logger, roomService)))
	mux.Handle("DELETE /rooms/{id}", requireAuth(delete.New(logger, roomService)))
}
