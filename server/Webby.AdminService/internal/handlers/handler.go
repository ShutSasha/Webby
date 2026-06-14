package handlers

import (
	"context"
	"webby/admin-service/internal/models"

	"github.com/google/uuid"
)

type service interface {
	ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error)
	AcceptComplaint(ctx context.Context, complaintID, userID uuid.UUID) error
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{service}
}
