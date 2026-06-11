package handlers

import (
	"context"
	"webby/vote-service/internal/services"

	"github.com/google/uuid"
)

type service interface {
	CreateWithRightChoice(ctx context.Context, roomID, userID uuid.UUID, voteText string, duration int, choices []string) error
	ListVotes(ctx context.Context, roomID, userID uuid.UUID) ([]services.VoteDetail, error)
	GetVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
	CastVote(ctx context.Context, voteId, choiceId, userId uuid.UUID) (*services.VoteDetail, error)
	RemoveVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
	DeleteVote(ctx context.Context, voteId, userId uuid.UUID) error
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{
		service: service,
	}
}
