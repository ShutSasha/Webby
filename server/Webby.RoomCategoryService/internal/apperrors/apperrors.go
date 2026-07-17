package apperrors

import "errors"

var (
	ErrInvalidInput          = errors.New("invalid input")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)
