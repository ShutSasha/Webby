package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Id        uuid.UUID
	RoomId    *uuid.UUID
	CreatedAt time.Time
}
