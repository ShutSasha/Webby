package apperrors

import "errors"

var (
	ErrPrivateChatAlreadyExists = errors.New("private chat already exists")
	ErrMessageNotFound          = errors.New("message not found")
	ErrInvalidInput             = errors.New("invalid input")
	ErrChatNotFound             = errors.New("chat not found")
	ErrNotFollowed              = errors.New("not followed to each other")
	ErrNotMemeber               = errors.New("not member of chat")
	ErrNotSender                = errors.New("not sender of message")
	ErrMemberNotFound           = errors.New("member not found")
)
