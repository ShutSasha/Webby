package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	SenderID  uuid.UUID
	ChatID    uuid.UUID
	Content   string
	IsEdited  bool
	EditedAt  *time.Time
	CreatedAt time.Time
}
