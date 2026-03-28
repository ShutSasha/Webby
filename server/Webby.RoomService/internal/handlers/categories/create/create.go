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

// createCategory godoc
// @Summary      Create a new category
// @Description  Create a new room category for organizing rooms.
// @Description  **Validation Rules:**
// @Description  * Name is **required** (cannot be null, empty, or whitespace-only)
// @Description  * Name length must be **2-50 characters**
// @Description  * Category name must be **unique** (409 Conflict if already exists)
// @Description  **Security Note:**
// @Description  * **ONLY ADMINISTRATORS** are authorized to perform this action.
// @Tags         RoomCategory
// @Accept       json
// @Produce      json
// @Param        request body docs.CreateCategoryRequest true "**Category data.** Name (2-50 chars, required, unique)"
// @Success      201 {object} docs.ApiResponse[docs.CategoryResponse] "Category successfully created"
// @Failure      400 {object} docs.Error400Response "Invalid input data: name missing, empty, too short (<2), too long (>50), or whitespace-only. Also: malformed JSON, wrong Content-Type"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Not authorized - Admin role required"
// @Failure      409 {object} docs.Error409Response "Conflict: Category name already exists"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /categories [post]
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
