package apperrors

import "errors"

var (
	ErrQueueItemNotFound = errors.New("queue item not found")
	ErrVideoNotFound     = errors.New("video not found")
	ErrConflict          = errors.New("queue item already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal server error")
)
