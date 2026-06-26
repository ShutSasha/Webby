package models

import (
	"time"

	"github.com/google/uuid"
)

type Complainer struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type Target struct {
	Type string    `json:"type"`
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Complaint struct {
	ID             uuid.UUID  `json:"id"`
	Complainer     Complainer `json:"complainer"`
	Target         Target     `json:"target"`
	ReasonType     string     `json:"reasonType"`
	AdditionalInfo *string    `json:"additionalInfo"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type ComplaintResult struct {
	ComplaintID uuid.UUID
	AdminID     uuid.UUID
	IsAccepted  bool
	CreatedAt   time.Time
}
