package models

import "github.com/google/uuid"

type RoomMember struct {
	RoomId     uuid.UUID
	UserId     uuid.UUID
	RoomPoints int
}

type MemberStatus struct {
	IsMember bool
	IsBanned bool
}
