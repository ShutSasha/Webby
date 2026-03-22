package errors

import (
	"errors"
	"log/slog"
	"net/http"
	"webby/internal/handlers/responses"
)

var (
	ErrNotFound = errors.New("not found")
)

type ValidationError map[string]string

func (e ValidationError) Error() string {
	return "validation failed"
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandler(logger *slog.Logger, h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {

			var valErr ValidationError
			if errors.As(err, &valErr) {
				logger.Warn("validation failed", slog.Any("problems", valErr))
				responses.Error(w, r, http.StatusBadRequest, "Validation failed", valErr)
				return
			}

			if errors.Is(err, ErrNotFound) {
				responses.Error(w, r, http.StatusNotFound, "Resourse is not found", map[string]string{
					"message": err.Error(),
				})
			} else {
				logger.Error("internal error", slog.String("error", err.Error()))
			}
		}
	}
}
