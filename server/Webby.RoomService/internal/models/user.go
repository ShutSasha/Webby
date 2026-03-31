package models

import "github.com/google/uuid"

type User struct {
	UserId    uuid.UUID
	Username  string
	AvatarUrl string
}

type RoomMemberInfo struct {
	UserId     uuid.UUID
	Username   string
	AvatarUrl  string
	RoomPoints int
}
