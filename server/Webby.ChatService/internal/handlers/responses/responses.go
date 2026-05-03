package responses

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"webby-chat/internal/apperrors"
	"webby-chat/pkg/http/render"
)

type ApiResponse[T any] struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *T                `json:"data"`
	Errors  map[string]string `json:"errors"`
}

type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Limit int `json:"pageSize"`
	Total int `json:"totalCount"`
}

func Error(w http.ResponseWriter, r *http.Request, statusCode int, message string, problems map[string]string) {
	render.Encode(w, r, statusCode, ApiResponse[struct{}]{
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

func HandleValidationError(c *gin.Context, err error) {
	problems := make(map[string]string)
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			problems[fe.Field()] = formatErrorMessage(fe)
		}
	} else {
		problems["message"] = err.Error()
	}

	c.JSON(http.StatusBadRequest, ApiResponse[struct{}]{
		Success: false,
		Message: "Validation error",
		Errors:  problems,
	})
}

func mapAppErrorToStatus(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, apperrors.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, apperrors.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func formatErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "notblank":
		return "this field cannot be empty or contain only spaces"
	case "email":
		return "invalid email format"
	case "uuid":
		return "invalid UUID format"
	case "url":
		return "invalid URL format"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters long", fe.Param())
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must not exceed %s characters", fe.Param())
		}
		return fmt.Sprintf("must not be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	default:
		return fmt.Sprintf("validation failed on the '%s' tag", fe.Tag())
	}
}

func HandleAppError(c *gin.Context, message string, err error) {
	status := mapAppErrorToStatus(err)
	c.JSON(status, ApiResponse[struct{}]{
		Success: false,
		Message: message,
		Errors:  map[string]string{"message": err.Error()},
	})
}
