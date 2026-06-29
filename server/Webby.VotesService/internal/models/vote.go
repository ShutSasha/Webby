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
	Status    string    `json:"status"`
}

type EnrichedVoting struct {
	ID        uuid.UUID `json:"id"`
	VoteText  string    `json:"voteText"`
	Duration  int       `json:"duration"`
	ExpiresAt time.Time `json:"expiresAt"`
	IsLocked  bool      `json:"isLocked"`
	MyVote    *string   `json:"myVote,omitempty"`
	Choices   []string  `json:"choices"`
}

type NextVideoInfo struct {
	Exists    bool       `json:"exists"`
	Duration  *int       `redis:"duration" json:"duration,omitempty"`
	ExpiresAt *time.Time `redis:"expires_at" json:"expiresAt,omitempty"`
}
