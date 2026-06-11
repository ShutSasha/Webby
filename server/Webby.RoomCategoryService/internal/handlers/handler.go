package handlers

import (
	"context"
)

type service interface {
	Create(ctx context.Context, name string) error
	List(ctx context.Context, search string, page, limit int) ([]string, int, error)
	Update(ctx context.Context, oldName string, newName string) error
	Delete(ctx context.Context, name string) error
}

type handler struct {
	service service
}

func New(service service) handler {
	return handler{service}
}
