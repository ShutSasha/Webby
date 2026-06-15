package models

import "github.com/google/uuid"

type Video struct {
	ID    uuid.UUID
	Title string
}
