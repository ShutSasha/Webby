package models

import "github.com/google/uuid"

type Reaction struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Cost       int       `json:"cost"`
	StickerURL string    `json:"stickerUrl"`
}
