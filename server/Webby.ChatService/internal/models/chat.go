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

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
}

type EnrichedChat struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type lastMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatHistoryItem struct {
	ChatID      uuid.UUID   `json:"chatId"`
	User        User        `json:"user"`
	LastMessage lastMessage `json:"lastMessage"`
}
