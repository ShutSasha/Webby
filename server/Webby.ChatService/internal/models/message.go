package models

import (
	"time"

	"github.com/google/uuid"
)

type Sender struct {
	ID        uuid.UUID
	Username  string
	AvatarURL string
}

type RichMessage struct {
	ID        uuid.UUID
	Sender    Sender
	ChatID    uuid.UUID
	Content   string
	IsEdited  bool
	EditedAt  *time.Time
	CreatedAt time.Time
}

type Message struct {
	ID        uuid.UUID
	SenderID  uuid.UUID
	ChatID    uuid.UUID
	Content   string
	IsEdited  bool
	EditedAt  *time.Time
	CreatedAt time.Time
}
