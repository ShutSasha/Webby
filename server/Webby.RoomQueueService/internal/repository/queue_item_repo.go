package repository

import (
	"context"
	"errors"
	"fmt"
	"webby-room-queue/internal/apperrors"
	"webby-room-queue/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueueItemRepository struct {
	db *pgxpool.Pool
}

func NewQueueItemRepository(db *pgxpool.Pool) *QueueItemRepository {
	return &QueueItemRepository{db: db}
}

func (r *QueueItemRepository) Create(
	ctx context.Context, item *models.QueueItem,
) (uuid.UUID, error) {
	const op = "repository.QueueItemRepository.Create"

	if item == nil {
		return uuid.Nil, fmt.Errorf(
			"%s: %w: item cannot be nil", op, apperrors.ErrInvalidInput,
		)
	}

	item.Id = uuid.New()

	query := `
		INSERT INTO queue_items (id, room_id, entity_id, entity_type, is_active, position)
		VALUES ($1, $2, $3, $4, $5,
			(SELECT COALESCE(MAX(position), 0) + 1
			 FROM queue_items WHERE room_id = $2))
		RETURNING created_at, position
	`

	err := r.db.QueryRow(
		ctx, query,
		item.Id, item.RoomId, item.EntityId, item.EntityType, item.IsActive,
	).Scan(&item.CreatedAt, &item.Position)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return item.Id, nil
}

func (r *QueueItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.QueueItemRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf(
			"%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput,
		)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var roomId uuid.UUID
	var position int
	err = tx.QueryRow(
		ctx,
		`SELECT room_id, position FROM queue_items WHERE id = $1`,
		id,
	).Scan(&roomId, &position)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrNotFound,
			)
		}
		return fmt.Errorf("%s: get item: %w", op, err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM queue_items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("%s: delete failed: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE queue_items SET position = position - 1
		 WHERE room_id = $1 AND position > $2`,
		roomId, position,
	)
	if err != nil {
		return fmt.Errorf("%s: shift positions: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (r *QueueItemRepository) GetById(
	ctx context.Context, id uuid.UUID,
) (*models.QueueItem, error) {
	const op = "repository.QueueItemRepository.GetById"

	if id == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, room_id, entity_id, entity_type,
		       is_active, position, created_at
		FROM queue_items
		WHERE id = $1
	`

	var item models.QueueItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.Id,
		&item.RoomId,
		&item.EntityId,
		&item.EntityType,
		&item.IsActive,
		&item.Position,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &item, nil
}

func (r *QueueItemRepository) ListByRoom(
	ctx context.Context, roomId uuid.UUID,
) ([]models.QueueItem, error) {
	const op = "repository.QueueItemRepository.ListByRoom"

	if roomId == uuid.Nil {
		return nil, fmt.Errorf(
			"%s: %w: invalid room id", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		SELECT id, room_id, entity_id, entity_type,
		       is_active, position, created_at
		FROM queue_items
		WHERE room_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.Query(ctx, query, roomId)
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

func (r *QueueItemRepository) MoveToTop(
	ctx context.Context, id uuid.UUID,
) error {
	const op = "repository.QueueItemRepository.MoveToTop"

	if id == uuid.Nil {
		return fmt.Errorf(
			"%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput,
		)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var roomId uuid.UUID
	var currentPos int
	err = tx.QueryRow(
		ctx,
		`SELECT room_id, position FROM queue_items WHERE id = $1`,
		id,
	).Scan(&roomId, &currentPos)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrNotFound,
			)
		}
		return fmt.Errorf("%s: get item: %w", op, err)
	}

	if currentPos == 1 {
		return nil
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE queue_items SET position = position + 1
		 WHERE room_id = $1 AND position < $2`,
		roomId, currentPos,
	)
	if err != nil {
		return fmt.Errorf("%s: shift positions: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE queue_items SET position = 1 WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("%s: update position: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}
