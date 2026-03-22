package services

import (
	"webby/internal/models"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(room *models.Room) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetById(id uuid.UUID) (*models.Room, error)
	ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error)
	ListPublic(page int, limit int) ([]models.Room, int64, error)
	Update(room *models.Room) (uuid.UUID, error)
}

type RoomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (r *RoomService) Create(room *models.Room) (uuid.UUID, error) {
	id, err := r.repo.Create(room)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *RoomService) Delete(id uuid.UUID) error {
	err := r.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoomService) GetById(id uuid.UUID) (*models.Room, error) {
	room, err := r.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomService) ListMy(userId uuid.UUID, page int, limit int) ([]models.Room, int64, error) {
	rooms, total, err := r.repo.ListMy(userId, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) ListPublic(page int, limit int) ([]models.Room, int64, error) {
	rooms, total, err := r.repo.ListPublic(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *RoomService) Update(room *models.Room) (uuid.UUID, error) {
	id, err := r.repo.Update(room)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
