package services

import (
	"context"
	"fmt"
	"webby/room-category-service/internal/apperrors"
	"webby/room-category-service/internal/models"
)

type Repository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, search string, offset int, limit int) ([]models.Category, int64, error)
	Update(ctx context.Context, oldName string, newName string) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (c *Service) Create(ctx context.Context, name string) error {
	const op = "services.CategoryService.Create"

	if name == "" {
		return fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	return c.repo.Create(ctx, name)
}

func (c *Service) Delete(ctx context.Context, name string) error {
	const op = "services.CategoryService.Delete"

	if name == "" {
		return fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	return c.repo.Delete(ctx, name)
}

func (c *Service) List(ctx context.Context, search string, page int, limit int) ([]models.Category, int64, error) {
	offset := (page - 1) * limit

	categories, total, err := c.repo.List(ctx, search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (c *Service) Update(ctx context.Context, oldName string, newName string) error {
	const op = "services.CategoryService.Update"

	if oldName == "" {
		return fmt.Errorf("%s: %w: old category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	if newName == "" {
		return fmt.Errorf("%s: %w: new category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	return c.repo.Update(ctx, oldName, newName)
}

func (c *Service) Exists(ctx context.Context, name string) (bool, error) {
	const op = "services.CategoryService.Exists"

	if name == "" {
		return false, fmt.Errorf("%s: %w: category name cannot be empty", op, apperrors.ErrInvalidInput)
	}

	return c.repo.Exists(ctx, name)
}
