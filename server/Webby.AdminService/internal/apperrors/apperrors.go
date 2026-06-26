package apperrors

import "errors"

var (
	ErrAlreadyResolved   = errors.New("complaint already resolved")
	ErrCantBanUser       = errors.New("can not ban user")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyBanned = errors.New("user is already banned")
	ErrVideoNotFound     = errors.New("video not found")
	ErrVideoAlreadyBanned     = errors.New("video is already banned")
)
