package services

import (
	"context"
	"fmt"
	"webby/room-category-service/internal/models"
)

type repository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, search string, offset, limit int) ([]models.Category, int, error)
	Update(ctx context.Context, oldName string, newName string) error
}

type service struct {
	repository repository
}

func New(repository repository) *service {
	return &service{repository}
}

func (c *service) Create(ctx context.Context, name string) error {
	const op = "service.Create"

	err := c.repository.Create(ctx, name)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to create category", op, err)
	}

	return nil
}

func (c *service) Delete(ctx context.Context, name string) error {
	const op = "service.Delete"

	err := c.repository.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to delete category", op, err)
	}

	return nil
}

func (c *service) List(ctx context.Context, search string, page, limit int) ([]string, int, error) {
	const op = "service.List"

	offset := (page - 1) * limit
	categories, total, err := c.repository.List(ctx, search, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: failed to list categories", op, err)
	}

	items := make([]string, len(categories))
	for i, cat := range categories {
		items[i] = cat.Name
	}

	return items, total, nil
}

func (c *service) Update(ctx context.Context, oldName string, newName string) error {
	const op = "service.Update"

	err := c.repository.Update(ctx, oldName, newName)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to update category", op, err)
	}

	return nil
}

func (c *service) Exists(ctx context.Context, name string) (bool, error) {
	const op = "service.Exists"

	exists, err := c.repository.Exists(ctx, name)
	if err != nil {
		return false, fmt.Errorf("%s: %w: failed to check if category exists", op, err)
	}

	return exists, nil
}
