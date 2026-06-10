package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"webby/chat-service/internal/apperrors"
	"webby/chat-service/pkg/http/render"
	"webby/chat-service/pkg/logger"

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

func Error(w http.ResponseWriter, r *http.Request, statusCode int, message string, problems map[string]string) {
	render.Encode(w, r, statusCode, ApiResponse[struct{}]{
		Success: false,
		Message: message,
		Errors:  problems,
	})
}

func HandleValidationError(c *gin.Context, err error) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	problems := make(map[string]string)
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		log.Debug("validation error", slog.String("err", err.Error()))
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

func mapAppErrorToStatus(err error) int {
	switch {
	case errors.Is(err, apperrors.ErrChatNotFound),
		errors.Is(err, apperrors.ErrMessageNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperrors.ErrPrivateChatAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, apperrors.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, apperrors.ErrNotMemeber),
		errors.Is(err, apperrors.ErrNotSender),
		errors.Is(err, apperrors.ErrNotFollowed):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func mapAppErrorToClientMessage(err error) string {
	switch {
	case errors.Is(err, apperrors.ErrChatNotFound):
		return "The requested chat was not found"
	case errors.Is(err, apperrors.ErrMessageNotFound):
		return "The requested message was not found"
	case errors.Is(err, apperrors.ErrPrivateChatAlreadyExists):
		return "A private chat with this user already exists"
	case errors.Is(err, apperrors.ErrInvalidInput):
		return "The provided input is invalid"
	case errors.Is(err, apperrors.ErrNotMemeber):
		return "User is not a member of this chat"
	case errors.Is(err, apperrors.ErrNotSender):
		return "You are not a sender of this message"
	case errors.Is(err, apperrors.ErrNotFollowed):
		return "You have to be followed to each other"
	default:
		return "An unexpected error occurred"
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
