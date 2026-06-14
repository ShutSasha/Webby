package handlers

import (
	"context"
	"webby/admin-service/internal/models"
)

type service interface {
	ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error)
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{service}
}
