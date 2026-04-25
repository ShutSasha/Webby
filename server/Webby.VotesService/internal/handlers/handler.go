package handlers

import (
	"context"
	"webby-vote-service/internal/services"

	"github.com/google/uuid"
)

type Service interface {
	CreateVote(
		ctx context.Context,
		roomId, userId uuid.UUID,
		voteType, voteText string,
		durationSeconds int,
		choices []services.CreateChoiceInput,
	) (*services.VoteDetail, error)
	ListVotes(ctx context.Context, roomId, userId uuid.UUID) ([]services.VoteDetail, error)
	GetVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
	CastVote(ctx context.Context, voteId, choiceId, userId uuid.UUID) (*services.VoteDetail, error)
	RemoveVote(ctx context.Context, voteId, userId uuid.UUID) (*services.VoteDetail, error)
	DeleteVote(ctx context.Context, voteId, userId uuid.UUID) error
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{
		service: service,
	}
}
