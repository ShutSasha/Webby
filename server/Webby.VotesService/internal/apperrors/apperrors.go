package apperrors

import "errors"

var (
	ErrNotHost        = errors.New("not a host")
	ErrAlreadyClosed  = errors.New("voting is already closed")
	ErrMemberNotFound = errors.New("member not found")
	ErrInvalidChoice  = errors.New("invalid choice")
	ErrVotingLocked   = errors.New("voting is locked")
)
