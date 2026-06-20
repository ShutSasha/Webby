package services

import (
	"context"
	"fmt"
)

type userStatsRetriever interface {
	GetMonthlyRegistrations(ctx context.Context) (map[string]int, error)
	GetMonthlySubscriptions(ctx context.Context) (map[string]int, error)
	GetTotalRegistrations(ctx context.Context) (int, error)
	GetMonthRevenue(ctx context.Context) (int, error)
}

type roomStatsRetriever interface {
	GetTotalRooms(ctx context.Context) (int, error)
}

type reportsStatsRetriever interface {
	GetActiveReportsCount(ctx context.Context) (int, error)
}

type statsService struct {
	userStatsRetriever    userStatsRetriever
	roomStatsRetriever    roomStatsRetriever
	reportsStatsRetriever reportsStatsRetriever
}

func NewStatsService(userStatsRetriever userStatsRetriever, roomStatsRetriever roomStatsRetriever, reportsStatsRetriever reportsStatsRetriever) *statsService {
	return &statsService{
		userStatsRetriever:    userStatsRetriever,
		roomStatsRetriever:    roomStatsRetriever,
		reportsStatsRetriever: reportsStatsRetriever,
	}
}

func (s *statsService) RegistrationsStats(ctx context.Context) (map[string]int, error) {
	const op = "statsSerice.RegistrationsStats"

	stats, err := s.userStatsRetriever.GetMonthlyRegistrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stats, nil
}

func (s *statsService) SubscriptionsStats(ctx context.Context) (map[string]int, error) {
	const op = "statsSerice.RegistrationsStats"

	stats, err := s.userStatsRetriever.GetMonthlySubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return stats, nil
}

func (s *statsService) GetAbsoluteStats(ctx context.Context) (map[string]int, error) {
	const op = "statsSerice.GetAbsoluteStats"

	absoluteStats := make(map[string]int)

	totalRegistrations, err := s.userStatsRetriever.GetTotalRegistrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	absoluteStats["totalRegistrations"] = totalRegistrations

	revenue, err := s.userStatsRetriever.GetMonthRevenue(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	absoluteStats["monthRevenue"] = revenue

	totalRooms, err := s.roomStatsRetriever.GetTotalRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	absoluteStats["totalRooms"] = totalRooms

	activeReportsCount, err := s.reportsStatsRetriever.GetActiveReportsCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	absoluteStats["activeReportsCount"] = activeReportsCount

	return absoluteStats, nil
}
