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

type user struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarUrl string    `json:"avatarUrl"`
}

type lastMessage struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatHistoryItem struct {
	ChatID      uuid.UUID   `json:"chatId"`
	User        user        `json:"user"`
	LastMessage lastMessage `json:"lastMessage"`
}
