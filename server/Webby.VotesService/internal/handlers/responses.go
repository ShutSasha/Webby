package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"webby/vote-service/internal/apperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiResponse[T any] struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *T                `json:"data"`
	Errors  map[string]string `json:"errors"`
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

func formatErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "uuid":
		return "must be a valid UUID"
	case "oneof":
		return fmt.Sprintf(
			"must be one of: %s", fe.Param(),
		)
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf(
				"must be at least %s characters long", fe.Param(),
			)
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf(
				"must not exceed %s characters", fe.Param(),
			)
		}
		return fmt.Sprintf("must not be greater than %s", fe.Param())
	default:
		return fmt.Sprintf(
			"validation failed on the '%s' tag", fe.Tag(),
		)
	}
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

func HandleAppError(c *gin.Context, message string, err error) {
	status := mapAppErrorToStatus(err)
	c.JSON(status, ApiResponse[struct{}]{
		Success: false,
		Message: message,
		Errors:  map[string]string{"message": err.Error()},
	})
}
