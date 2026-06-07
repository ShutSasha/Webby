package services

import (
	"context"
	"fmt"
)

type presenceRepo interface {
	AddUser(ctx context.Context, chatID, userID string) error
	RemoveUser(ctx context.Context, chatID, userID string) error
}

type presenceService struct {
	presenceRepo presenceRepo
}

func NewPresenceService(presenceRepo presenceRepo) *presenceService {
	return &presenceService{presenceRepo}
}

func (p *presenceService) AddUser(ctx context.Context, chatID, userID string) error {
	const op = "services.presenceService.AddUser"

	err := p.presenceRepo.AddUser(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *presenceService) RemoveUser(ctx context.Context, chatID, userID string) error {
	const op = "services.presenceService.RemoveUser"

	err := p.presenceRepo.RemoveUser(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
