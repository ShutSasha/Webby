package handlers

import (
	"context"
	"webby/vote-service/internal/models"

	"github.com/google/uuid"
)

type service interface {
	CreateWithRightChoice(ctx context.Context, roomID, userID uuid.UUID, voteText string, duration int, choices []string) error
	ResolveVoting(ctx context.Context, roomID, userID, voteID uuid.UUID, rightChoice string) error
	ListVotings(ctx context.Context, roomID, userID uuid.UUID) ([]models.EnrichedVoting, error)
	CastVote(ctx context.Context, voteId, choiceId, userId uuid.UUID) error
	RemoveVote(ctx context.Context, voteId, userId uuid.UUID) error
	DeleteVote(ctx context.Context, voteId, userId uuid.UUID) error
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{service}
}
