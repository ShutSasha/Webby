package models

import (
	"time"

	"github.com/google/uuid"
)

type QueueItem struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	VideoID   string
	IsActive  bool
	Position  int
	CreatedAt time.Time
}

type EnrichedQueueItem struct {
	ID        uuid.UUID
	VideoID   string
	Title     string
	Thumbnail string
	VideoUrl  string
	IsActive  bool
	Position  int
}
