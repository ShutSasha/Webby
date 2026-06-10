package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func HandleValidationError(c *gin.Context, err error) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	problems := make(map[string]string)
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			problems[fe.Field()] = formatErrorMessage(fe)
		}
	} else {
		log.Debug("non-validation error in request binding", slog.String("err", err.Error()))
		problems["message"] = "invalid request body"
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
	case errors.Is(err, apperrors.ErrQueueItemNotFound),
		errors.Is(err, apperrors.ErrVideoNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, apperrors.ErrNotRoomMember):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func mapAppErrorToClientMessage(err error) string {
	switch {
	case errors.Is(err, apperrors.ErrQueueItemNotFound):
		return "The requested queue item was not found"
	case errors.Is(err, apperrors.ErrVideoNotFound):
		return "The requested video was not found"
	case errors.Is(err, apperrors.ErrConflict):
		return "This video already exists in the queue"
	case errors.Is(err, apperrors.ErrNotRoomMember):
		return "You are not a room member"
	default:
		return "An unexpected error occurred"
	}
}

func HandleAppError(c *gin.Context, message string, err error) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	log.Error(message, slog.String("error", err.Error()))

	status := mapAppErrorToStatus(err)
	clientMessage := mapAppErrorToClientMessage(err)

	c.JSON(status, ApiResponse[struct{}]{
		Success: false,
		Message: message,
		Errors:  map[string]string{"message": clientMessage},
	})
}
