package services

import (
	"context"
	"fmt"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
)

type repository interface {
	ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error)
}

type notificationSender interface {
	SendNotificationToUser(ctx context.Context, userID, targetID uuid.UUID) error
}

type service struct {
	repository         repository
	notificationSender notificationSender
}

func New(repository repository, notificationSender notificationSender) *service {
	return &service{
		repository:         repository,
		notificationSender: notificationSender,
	}
}

func (s *service) ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error) {
	const op = "service.ListComplaints"

	offset := (page - 1) * limit
	result, total, err := s.repository.ListComplaints(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return result, total, nil
}

func (s *service) AcceptComplaint(ctx context.Context, complaintID, userID uuid.UUID) error {
	const op = "service.AcceptComplaint"

	

	err := s.repository.AcceptComplaint(ctx, complaintID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.notificationSender.SendNotificationToUser(ctx)

	return nil
}
