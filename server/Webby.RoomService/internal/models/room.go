package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id           uuid.UUID
	HostId       uuid.UUID
	CategoryName string
	Name         string
	Thumbnail    string
	IsPrivate    bool
	CreatedAt    time.Time
	ChatId       *uuid.UUID // populated from ChatService, not stored in DB
}

type PublicRoom struct {
	Id            uuid.UUID
	HostId        uuid.UUID
	HostUsername  string
	HostAvatarUrl string
	CategoryName  string
	Name          string
	Thumbnail     string
	IsPrivate     bool
	CreatedAt     time.Time
}
