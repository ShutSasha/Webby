package render

import (
	"github.com/goccy/go-json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"webby/internal/apperrors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	_ = validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func DecodeValid[T any](r *http.Request) (T, map[string]string, error) {
	var v T
	problems := make(map[string]string)

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("%w: %w", apperrors.ErrInvalidInput, err)
	}

	if err := validate.Struct(v); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				problems[fieldError.Field()] = formatErrorMessage(fieldError)
			}
		} else {
			return v, nil, fmt.Errorf("%w: %w", apperrors.ErrInvalidInput, err)
		}
	}

	if len(problems) > 0 {
		return v, problems, fmt.Errorf("%w: %d problems", apperrors.ErrInvalidInput, len(problems))
	}

	return v, nil, nil
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
