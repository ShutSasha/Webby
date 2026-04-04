package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID
	SenderId  uuid.UUID
	ChatId    uuid.UUID
	Content   string
	IsEdited  bool
	EditedAt  *time.Time
	CreatedAt time.Time
}
