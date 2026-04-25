package handlers

import (
	"context"
	"webby/internal/models"
)

type Service interface {
	Create(ctx context.Context, name string) error
	List(ctx context.Context, search string, page int, limit int) ([]models.Category, int64, error)
	Update(ctx context.Context, oldName string, newName string) error
	Delete(ctx context.Context, name string) error
}

type handler struct {
	service Service
}

func New(service Service) handler {
	return handler{
		service: service,
	}
}
