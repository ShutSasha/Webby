package apperrors

import "errors"

var (
	ErrNotHost            = errors.New("not host")
	ErrNotMember          = errors.New("not member")
	ErrRemoveHost         = errors.New("only host can perform")
	ErrRoomNotFound       = errors.New("room not found")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrRoomMemberNotFound = errors.New("room member not found")
	ErrBanned             = errors.New("banned")
	ErrMaxMembersReached  = errors.New("max members reached")
)
