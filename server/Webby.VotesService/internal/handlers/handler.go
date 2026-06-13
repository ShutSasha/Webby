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
	CastVote(ctx context.Context, roomID, voteID, userID uuid.UUID, choice string) error
	CreateVotingForNextVideo(ctx context.Context, roomID, userID uuid.UUID) error
	VoteForNextVideo(ctx context.Context, roomID, userID, queueItemID uuid.UUID) error
	HasNextVideoVoting(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{service}
}
