package models

import "github.com/google/uuid"

type ChatMember struct {
	ChatId uuid.UUID
	UserId uuid.UUID
}
