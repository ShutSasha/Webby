package apperrors

import "errors"

var (
	ErrQueueItemNotFound = errors.New("queue item not found")
	ErrVideoNotFound     = errors.New("videow not found")
	ErrConflict          = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal server error")
)
