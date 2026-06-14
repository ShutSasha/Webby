package services

import (
	"context"
	"fmt"
	"webby/admin-service/internal/models"
)

type repository interface {
	ListComplaints(ctx context.Context, offset, limit int) ([]models.Complaint, int, error)
}

type service struct {
	repository repository
}

func New(repository repository) *service {
	return &service{repository}
}

func (c *service) ListComplaints(ctx context.Context, page, limit int) ([]models.Complaint, int, error) {
	const op = "service.ListComplaints"

	offset := (page - 1) * limit
	result, total, err := c.repository.ListComplaints(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return result, total, nil
}
