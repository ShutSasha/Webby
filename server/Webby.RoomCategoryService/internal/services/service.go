package services

import (
	"context"
	"fmt"
	"webby/room-category-service/internal/apperrors"
	"webby/room-category-service/internal/models"
)

type repository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, search string, offset int, limit int) ([]models.Category, int64, error)
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

	if name == "" {
		return fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	err := c.repository.Create(ctx, name)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to create category", op, err)
	}

	return nil
}

func (c *service) Delete(ctx context.Context, name string) error {
	const op = "service.Delete"

	if name == "" {
		return fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	err := c.repository.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to delete category", op, err)
	}

	return nil
}

func (c *service) List(ctx context.Context, search string, page int, limit int) ([]models.Category, int64, error) {
	const op = "service.List"

	offset := (page - 1) * limit
	categories, total, err := c.repository.List(ctx, search, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w: failed to list categories", op, err)
	}

	return categories, total, nil
}

func (c *service) Update(ctx context.Context, oldName string, newName string) error {
	const op = "service.Update"

	if oldName == "" {
		return fmt.Errorf("%s: %w: old category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	if newName == "" {
		return fmt.Errorf("%s: %w: new category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	err := c.repository.Update(ctx, oldName, newName)
	if err != nil {
		return fmt.Errorf("%s: %w: failed to update category", op, err)
	}

	return nil
}

func (c *service) Exists(ctx context.Context, name string) (bool, error) {
	const op = "service.Exists"

	if name == "" {
		return false, fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	exists, err := c.repository.Exists(ctx, name)
	if err != nil {
		return false, fmt.Errorf("%s: %w: failed to check if category exists", op, err)
	}

	return exists, nil
}
