package create

import (
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	httpErrors "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type Creator interface {
	Create(name string) (uuid.UUID, error)
}

func New(logger *slog.Logger, creator Creator) http.Handler {
	return httpErrors.MakeHandler(logger, createCategory(logger, creator))
}

// @Title Create a new category
// @Description Create a new room category. Name is required (2-50 chars, unique). Only administrators can perform this action.
// @Param  request  body  docs.CreateCategoryRequest  true  "Category data. Name (2-50 chars, required, unique)"
// @Success  201  object docs.CategoryApiResponse  "Category successfully created"
// @Failure  400  object docs.ErrorResponse  "Invalid input data"
// @Failure  401  object docs.ErrorResponse  "Missing or invalid authentication token"
// @Failure  403  object docs.ErrorResponse  "Not authorized - Admin role required"
// @Failure  409  object docs.ErrorResponse  "Conflict: Category name already exists"
// @Failure  500  object docs.ErrorResponse  "Internal server error"
// @Resource RoomCategory
// @Route /api/categories [post]
func createCategory(logger *slog.Logger, creator Creator) httpErrors.APIFunc {
	_ = logger.With(slog.String("operation", "httpserver.categories.create"))

	type request struct {
		Name string `json:"name" validate:"required,notblank,min=2,max=50"`
	}

	type response struct {
		Id   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		var req request
		req, problems, err := render.DecodeValid[request](r)
		if len(problems) > 0 {
			logger.Debug("problems debug", slog.Any("problems", problems))
			return responses.NewValidationError("Validation error", problems)
		}
		if err != nil {
			return responses.NewApiError("Validation error", err)
		}

		id, err := creator.Create(req.Name)
		if err != nil {
			return responses.NewApiError("Create category error", err)
		}

		render.Encode(w, r, http.StatusCreated, responses.ApiResponse[response]{
			Success: true,
			Message: "Category created",
			Data: &response{
				Id:   id,
				Name: req.Name,
			},
		})
		return nil
	}
}
