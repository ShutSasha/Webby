package services

import (
	"context"
	"fmt"
)

type registrationsRetriever interface {
	GetMonthlyRegistrations(ctx context.Context) (map[string]int, error)
	GetMonthlySubscriptions(ctx context.Context) (map[string]int, error)
}

type statsService struct {
	registrationsRetriever registrationsRetriever
}

func NewStatsService(registrationsRetriever registrationsRetriever) *statsService {
	return &statsService{
		registrationsRetriever: registrationsRetriever,
	}
}

func (s *statsService) RegistrationsStats(ctx context.Context) (map[string]int, error) {
	const op = "statsSerice.RegistrationsStats"

	stats, err := s.registrationsRetriever.GetMonthlyRegistrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stats, nil
}

func (s *statsService) SubscriptionsStats(ctx context.Context) (map[string]int, error) {
	const op = "statsSerice.RegistrationsStats"

	stats, err := s.registrationsRetriever.GetMonthlySubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stats, nil
}
