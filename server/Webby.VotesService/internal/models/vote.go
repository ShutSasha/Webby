package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	VoteText  string
	Duration  int
	CreatedAt time.Time
}

type EnrichedVoting struct {
	ID        uuid.UUID `json:"id"`
	VoteText  string    `json:"voteText"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
	Choices   []string  `json:"choices"`
}
