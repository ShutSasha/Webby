package apperrors

import "errors"

var (
	ErrConflict          = errors.New("queue item already exists")
	ErrNotRoomMember     = errors.New("not room member")
	ErrQueueItemNotFound = errors.New("queue item not found")
	ErrVideoNotFound     = errors.New("video not found")
)
