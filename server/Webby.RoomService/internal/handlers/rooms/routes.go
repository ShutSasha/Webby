package rooms

import (
	"log/slog"
	"net/http"
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/get"
	getByToken "webby/internal/handlers/rooms/get_by_token"
	listMembers "webby/internal/handlers/rooms/list_members"
	listMy "webby/internal/handlers/rooms/list_my"
	listPublic "webby/internal/handlers/rooms/list_public"
	removeMember "webby/internal/handlers/rooms/remove_member"
	"webby/internal/handlers/rooms/update"
	updatePoints "webby/internal/handlers/rooms/update_points"
	"webby/pkg/http/middleware/auth"
)

type RoomService interface {
	create.Creator
	listMy.MyLister
	listPublic.PublicLister
	get.Getter
	getByToken.TokenGetter
	update.Updater
	delete.Deleter
	listMembers.MemberLister
	removeMember.MemberRemover
	updatePoints.PointsUpdater
}

func RegisterRooms(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger,
	roomService RoomService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("GET /api/rooms/public", listPublic.New(logger, roomService))

	mux.Handle("POST /api/rooms", requireAuth(create.New(logger, roomService)))
	mux.Handle("POST /api/rooms/token", requireAuth(getByToken.New(logger, roomService)))
	mux.Handle("GET /api/rooms/my", requireAuth(listMy.New(logger, roomService)))
	mux.Handle("GET /api/rooms/{id}", requireAuth(get.New(logger, roomService)))
	mux.Handle("GET /api/rooms/{id}/members", requireAuth(listMembers.New(logger, roomService)))
	mux.Handle("DELETE /api/rooms/{id}/members/{memberId}", requireAuth(removeMember.New(logger, roomService)))
	mux.Handle("PATCH /api/rooms/{id}/members/{memberId}/points", requireAuth(updatePoints.New(logger, roomService)))
	mux.Handle("PUT /api/rooms/{id}", requireAuth(update.New(logger, roomService)))
	mux.Handle("DELETE /api/rooms/{id}", requireAuth(delete.New(logger, roomService)))
}
