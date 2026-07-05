package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID
	HostID    uuid.UUID
	Category  string
	Name      string
	Thumbnail string
	IsPrivate bool
	CreatedAt time.Time
	ChatID    *uuid.UUID
}

type PublicRoom struct {
	ID            uuid.UUID
	HostID        uuid.UUID
	HostUsername  string
	HostAvatarUrl string
	Category      string
	Name          string
	Thumbnail     string
	IsPrivate     bool
	CreatedAt     time.Time
}

type ActiveRoom struct {
	ID      uuid.UUID
	ChatID  string
	UserIDs []uuid.UUID
}
