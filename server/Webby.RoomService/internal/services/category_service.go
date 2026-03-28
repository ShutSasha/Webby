package services

import (
	"webby/internal/models"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(name string) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	List(search string, offset int, limit int) ([]models.Category, int64, error)
	Update(id uuid.UUID, name string) error
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (c *CategoryService) Create(name string) (uuid.UUID, error) {
	id, err := c.repo.Create(name)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (c *CategoryService) Delete(id uuid.UUID) error {
	err := c.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryService) List(search string, page int, limit int) ([]models.Category, int64, error) {
	offset := (page - 1) * limit
	
	categories, total, err := c.repo.List(search, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (c *CategoryService) Update(id uuid.UUID, name string) error {
	err := c.repo.Update(id, name)
	if err != nil {
		return err
	}
	return nil
}
