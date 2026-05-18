package models

import "github.com/google/uuid"

type User struct {
	UserID    uuid.UUID
	Username  string
	AvatarUrl string
}

type RoomMemberInfo struct {
	UserID     uuid.UUID
	Username   string
	AvatarUrl  string
	RoomPoints int
}
