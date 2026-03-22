package responses

import (
	"net/http"
	"webby/pkg/http/render"
)

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, statusCode int, message string, problems map[string]string) {
	render.Encode(w, r, statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  problems,
	})
}

type ApiError struct {
	Title string
	Err   error
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

func NewApiError(title string, err error) error {
	return &ApiError{Title: title, Err: err}
}

type ValidationError struct {
	Title string
	Err   map[string]string
}

func (e *ValidationError) Error() string {
	return e.Title
}

func NewValidationError(title string, problems map[string]string) error {
	return &ValidationError{Title: title, Err: problems}
}
