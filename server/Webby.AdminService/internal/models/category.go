package models

import (
	"time"

	"github.com/google/uuid"
)

type Complaint struct {
	ID             uuid.UUID `json:"id"`
	AuthorID       uuid.UUID `json:"authorId"`
	TargetType     string    `json:"targetType"`
	TargetID       uuid.UUID `json:"targetId"`
	ReasonType     string    `json:"reasonType"`
	AdditionalInfo string    `json:"additionalInfo"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ComplaintResult struct {
	ComplaintID uuid.UUID
	AdminID     uuid.UUID
	IsAccepted  bool
	CreatedAt   time.Time
}
