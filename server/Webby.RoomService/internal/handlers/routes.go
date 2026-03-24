package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"webby/internal/config"
	"webby/internal/handlers/categories"
	catCreate "webby/internal/handlers/categories/create"
	catDelete "webby/internal/handlers/categories/delete"
	catList "webby/internal/handlers/categories/list"
	catUpdate "webby/internal/handlers/categories/update"
	"webby/internal/handlers/rooms"
	"webby/internal/handlers/rooms/create"
	"webby/internal/handlers/rooms/delete"
	"webby/internal/handlers/rooms/get"
	getByToken "webby/internal/handlers/rooms/get_by_token"
	listMy "webby/internal/handlers/rooms/list_my"
	listPublic "webby/internal/handlers/rooms/list_public"
	"webby/internal/handlers/rooms/update"

	httpSwagger "github.com/swaggo/http-swagger"
)

type RoomService interface {
	create.Creator
	listMy.MyLister
	listPublic.PublicLister
	get.Getter
	getByToken.TokenGetter
	update.Updater
	delete.Deleter
}

type CategoryService interface {
	catCreate.Creator
	catList.Lister
	catUpdate.Updater
	catDelete.Deleter
}

func addRoutes(
	mux *http.ServeMux,
	cfg *config.Config,
	logger *slog.Logger,
	roomService RoomService,
	categoryService CategoryService,
) {
	mux.Handle("/api/", http.NotFoundHandler())

	rooms.RegisterRooms(mux, []byte(cfg.JwtSecret), logger, roomService)
	categories.RegisterCategories(mux, []byte(cfg.JwtSecret), logger, categoryService)

	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", cfg.Http.Host, cfg.Http.Port)),
		httpSwagger.DocExpansion("false"),
	))
	mux.Handle("GET /api/docs", http.RedirectHandler("/swagger/", http.StatusMovedPermanently))
}
