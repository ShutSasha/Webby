package services

import (
	"context"
	"fmt"
	"webby/room-service/internal/models"
)

type reactionRepository interface {
	List(ctx context.Context) ([]models.Reaction, error)
}

type reactionService struct {
	reactionRepo reactionRepository
}

func NewReactionService(reactionRepo reactionRepository) *reactionService {
	return &reactionService{reactionRepo: reactionRepo}
}

func (s *reactionService) List(ctx context.Context) ([]models.Reaction, error) {
	const op = "services.reactionService.List"

	reactions, err := s.reactionRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reactions, nil
}
