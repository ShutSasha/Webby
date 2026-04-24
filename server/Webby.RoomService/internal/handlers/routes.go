package httpserver

import (
	"log/slog"
	"net/http"
	"webby/internal/config"
	"webby/internal/handlers/roomqueues"
	"webby/internal/handlers/rooms"
	addMember "webby/internal/handlers/rooms/add_member"
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/get"
	listMembers "webby/internal/handlers/rooms/list_members"
	listMy "webby/internal/handlers/rooms/list_my"
	listPublic "webby/internal/handlers/rooms/list_public"
	removeMember "webby/internal/handlers/rooms/remove_member"
	"webby/internal/handlers/rooms/update"
	updatePoints "webby/internal/handlers/rooms/update_points"
	"webby/internal/handlers/votes"
)

type RoomService interface {
	create.Creator
	listMy.MyLister
	listPublic.PublicLister
	get.Getter
	update.Updater
	delete.Deleter
	listMembers.MemberLister
	addMember.MemberAdder
	removeMember.MemberRemover
	updatePoints.PointsUpdater
}

type QueueItemService interface {
	roomqueues.QueueItemService
}

type VoteService interface {
	votes.VoteService
}

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	logger *slog.Logger,
	roomService RoomService,
	queueItemService QueueItemService,
	voteService VoteService,
) {
	mux.Handle("/api/", http.NotFoundHandler())

	rooms.RegisterRooms(mux, []byte(cfg.JwtSecret), logger, roomService)
	roomqueues.RegisterRoomQueue(mux, []byte(cfg.JwtSecret), logger, queueItemService)
	votes.RegisterVotes(mux, []byte(cfg.JwtSecret), logger, voteService)

	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/oas.json")
	})

	mux.HandleFunc("GET /swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "./docs/oas.yml")
	})
}
