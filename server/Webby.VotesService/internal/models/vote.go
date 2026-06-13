package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID        uuid.UUID `redis:"id"`
	RoomID    uuid.UUID `redis:"room_id"`
	VoteText  string    `redis:"vote_text"`
	Duration  int       `redis:"duration"`
	CreatedAt time.Time `redis:"created_at"`
}

type EnrichedVoting struct {
	ID        uuid.UUID `json:"id"`
	VoteText  string    `json:"voteText"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"createdAt"`
	Choices   []string  `json:"choices"`
}
