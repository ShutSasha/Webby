package apperrors

import "errors"

var (
	ErrNotHost       = errors.New("not a host")
	ErrAlreadyClosed = errors.New("voting is already closed")
	ErrNotMember     = errors.New("not a member")
	ErrInvalidChoice = errors.New("invalid choice")
	ErrVotingLocked  = errors.New("voting is locked")
)
