package httpserver

import (
	"log/slog"
	"net/http"
	"webby/internal/config"
	"webby/internal/handlers/categories"
	catCreate "webby/internal/handlers/categories/create"
	catDelete "webby/internal/handlers/categories/delete"
	catList "webby/internal/handlers/categories/list"
	catUpdate "webby/internal/handlers/categories/update"
	"webby/internal/handlers/roomqueues"
	"webby/internal/handlers/rooms"
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
	"webby/internal/handlers/votes"
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

type CategoryService interface {
	catCreate.Creator
	catList.Lister
	catUpdate.Updater
	catDelete.Deleter
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
	categoryService CategoryService,
	queueItemService QueueItemService,
	voteService VoteService,
) {
	mux.Handle("/api/", http.NotFoundHandler())

	rooms.RegisterRooms(mux, []byte(cfg.JwtSecret), logger, roomService)
	categories.RegisterCategories(mux, []byte(cfg.JwtSecret), logger, categoryService)
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
