package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id         uuid.UUID
	HostId     uuid.UUID
	CategoryId uuid.UUID
	Name       string
	Thumbnail  string
	IsPrivate  bool
	CreatedAt  time.Time
}
