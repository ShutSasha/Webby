package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	Id              uuid.UUID
	RoomId          uuid.UUID
	Type            string
	VoteText        string
	CreatedAt       time.Time
	DurationSeconds int
}

type VoteChoice struct {
	Id          uuid.UUID
	VoteId      uuid.UUID
	Name        string
	Votes       int
	IsCorrect   bool
	QueueItemId *uuid.UUID
}

type UserVote struct {
	VoteChoiceId uuid.UUID
	UserId       uuid.UUID
}
