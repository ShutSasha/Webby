package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrForbidden      = errors.New("forbidden")
	ErrNotHost        = errors.New("not a host")
	ErrInternal       = errors.New("internal server error")
	ErrAlreadyClosed  = errors.New("voting is already closed")
	ErrMemberNotFound = errors.New("member not found")
	ErrInvalidChoice  = errors.New("invalid choice")
)
