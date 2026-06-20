package handlers

import (
	"context"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
)

type complaintsService interface {
	ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error)
	AcceptComplaint(ctx context.Context, complaintID, userID uuid.UUID) error
	DenyComplaint(ctx context.Context, complaintID, userID uuid.UUID, reason string) error
}

type statsService interface {
	RegistrationsStats(ctx context.Context) (map[string]int, error)
	SubscriptionsStats(ctx context.Context) (map[string]int, error)
}

type handler struct {
	complaintsService complaintsService
	statsService      statsService
}

func New(complaintsService complaintsService, statsService statsService) handler {
	return handler{
		complaintsService: complaintsService,
		statsService:      statsService,
	}
}
