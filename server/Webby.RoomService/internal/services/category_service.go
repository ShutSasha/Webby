package services

import (
	"context"
	"webby/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, search string, offset int, limit int) ([]models.Category, int64, error)
	Update(ctx context.Context, oldName string, newName string) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (c *CategoryService) Create(ctx context.Context, name string) error {
	return c.repo.Create(ctx, name)
}

func (c *CategoryService) Delete(ctx context.Context, name string) error {
	return c.repo.Delete(ctx, name)
}

func (c *CategoryService) List(ctx context.Context, search string, page int, limit int) ([]models.Category, int64, error) {
	offset := (page - 1) * limit

	categories, total, err := c.repo.List(ctx, search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (c *CategoryService) Update(ctx context.Context, oldName string, newName string) error {
	return c.repo.Update(ctx, oldName, newName)
}
