package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type QueueItemRepository struct {
	db *sql.DB
}

func NewQueueItemRepository(db *sql.DB) *QueueItemRepository {
	return &QueueItemRepository{db: db}
}

func (r *QueueItemRepository) Create(item *models.QueueItem) (uuid.UUID, error) {
	const op = "repository.QueueItemRepository.Create"

	if item == nil {
		return uuid.Nil, fmt.Errorf("%s: %w: item cannot be nil", op, apperrors.ErrInvalidInput)
	}

	item.Id = uuid.New()

	query := `
		INSERT INTO queue_items (id, room_id, entity_id, entity_type, is_active, position)
		VALUES ($1, $2, $3, $4, $5, (SELECT COALESCE(MAX(position), 0) + 1 FROM queue_items WHERE room_id = $2))
		RETURNING created_at, position
	`

	err := r.db.QueryRow(query, item.Id, item.RoomId, item.EntityId, item.EntityType, item.IsActive).Scan(&item.CreatedAt, &item.Position)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return item.Id, nil
}

func (r *QueueItemRepository) Delete(id uuid.UUID) error {
	const op = "repository.QueueItemRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM queue_items WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected failed: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: queue item %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}

func (r *QueueItemRepository) GetById(id uuid.UUID) (*models.QueueItem, error) {
	const op = "repository.QueueItemRepository.GetById"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, entity_id, entity_type, is_active, position, created_at
		FROM queue_items
		WHERE id = $1
	`

	var item models.QueueItem
	err := r.db.QueryRow(query, id).Scan(
		&item.Id,
		&item.RoomId,
		&item.EntityId,
		&item.EntityType,
		&item.IsActive,
		&item.Position,
		&item.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: queue item %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &item, nil
}

func (r *QueueItemRepository) ListByRoom(roomId uuid.UUID) ([]models.QueueItem, error) {
	const op = "repository.QueueItemRepository.ListByRoom"

	if roomId == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, entity_id, entity_type, is_active, position, created_at
		FROM queue_items
		WHERE room_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.Query(query, roomId)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	items := []models.QueueItem{}
	for rows.Next() {
		var item models.QueueItem
		err := rows.Scan(
			&item.Id,
			&item.RoomId,
			&item.EntityId,
			&item.EntityType,
			&item.IsActive,
			&item.Position,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return items, nil
}

func (r *QueueItemRepository) MoveToTop(id uuid.UUID) error {
	const op = "repository.QueueItemRepository.MoveToTop"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	var roomId uuid.UUID
	err = tx.QueryRow(`SELECT room_id FROM queue_items WHERE id = $1`, id).Scan(&roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: queue item %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return fmt.Errorf("%s: get room_id: %w", op, err)
	}

	var minPos int
	err = tx.QueryRow(`SELECT COALESCE(MIN(position), 1) FROM queue_items WHERE room_id = $1`, roomId).Scan(&minPos)
	if err != nil {
		return fmt.Errorf("%s: get min position: %w", op, err)
	}

	newPos := minPos - 1
	_, err = tx.Exec(`UPDATE queue_items SET position = $1 WHERE id = $2`, newPos, id)
	if err != nil {
		return fmt.Errorf("%s: update position: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}
