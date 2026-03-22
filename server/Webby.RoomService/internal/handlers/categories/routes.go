package categories

import (
	"log/slog"
	"net/http"
	"webby/internal/handlers/categories/create"
	"webby/internal/handlers/categories/delete"
	"webby/internal/handlers/categories/list"
	"webby/internal/handlers/categories/update"
	"webby/pkg/http/middleware/auth"
)

type CategoryService interface {
	create.Creator
	list.Lister
	update.Updater
	delete.Deleter
}

func RegisterCategories(mux *http.ServeMux, jwtSecret []byte, logger *slog.Logger,
	categoryService CategoryService) {
	requireAuth := auth.AuthMiddleware(jwtSecret)

	mux.Handle("GET /categories", list.New(logger, categoryService))

	mux.Handle("POST /categories", requireAuth(create.New(logger, categoryService)))
	mux.Handle("PUT /categories/{id}", requireAuth(update.New(logger, categoryService)))
	mux.Handle("DELETE /categories/{id}", requireAuth(delete.New(logger, categoryService)))
}
