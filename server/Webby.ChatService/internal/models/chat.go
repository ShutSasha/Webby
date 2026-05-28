package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	RoomID    *uuid.UUID
	CreatedAt time.Time
}
