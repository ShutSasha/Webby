package errors

import (
	"errors"
	"log/slog"
	"net/http"
	"webby-chat/internal/apperrors"
	"webby-chat/internal/handlers/responses"
	"webby-chat/pkg/logger"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log := logger.FromContext(r.Context())

			var apiErr *responses.ApiError
			displayTitle := "An error occurred"
			if errors.As(err, &apiErr) {
				displayTitle = apiErr.Title
			}

			unwrappedErr := err
			for errors.Unwrap(unwrappedErr) != nil {
				unwrappedErr = errors.Unwrap(unwrappedErr)
			}

			var valErr *responses.ValidationError
			if errors.As(err, &valErr) {
				log.Warn("validation failed", slog.Any("problems", valErr))
				responses.Error(w, r, http.StatusBadRequest, valErr.Title, valErr.Err)
				return
			}

			statusCode := http.StatusInternalServerError

			switch {
			case errors.Is(err, apperrors.ErrConflict):
				statusCode = http.StatusConflict
			case errors.Is(err, apperrors.ErrNotFound):
				statusCode = http.StatusNotFound
			case errors.Is(err, apperrors.ErrInvalidInput):
				statusCode = http.StatusBadRequest
			case errors.Is(err, apperrors.ErrForbidden):
				statusCode = http.StatusForbidden
			default:
				log.Error("internal server error", slog.String("error", err.Error()))
				responses.Error(w, r, statusCode, "Internal server error", map[string]string{
					"message": "something went wrong",
				})
				return
			}

			responses.Error(w, r, statusCode, displayTitle, map[string]string{
				"message": unwrappedErr.Error(),
			})
		}
	}
}
