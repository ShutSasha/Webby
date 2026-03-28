package models

import (
	"time"

	"github.com/google/uuid"
)

type QueueItem struct {
	Id         uuid.UUID
	RoomId     uuid.UUID
	EntityId   uuid.UUID
	EntityType string
	IsActive   bool
	CreatedAt  time.Time
}
