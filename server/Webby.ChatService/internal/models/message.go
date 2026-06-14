package models

import (
	"time"

	"github.com/google/uuid"
)

type Sender struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarUrl"`
}

type RichMessage struct {
	ID        uuid.UUID `json:"id"`
	Sender    Sender    `json:"sender"`
	Content   string    `json:"content"`
	IsEdited  bool      `json:"isEdited"`
	CreatedAt time.Time `json:"createdAt"`
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
