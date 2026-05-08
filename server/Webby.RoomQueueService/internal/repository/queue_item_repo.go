package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/internal/models"

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
) (uuid.UUID, int, error) {
	const op = "repository.QueueItemRepository.Create"

	if item == nil {
		return uuid.Nil, -1, fmt.Errorf(
			"%s: %w: item cannot be nil", op, apperrors.ErrInvalidInput,
		)
	}

	query := `
		INSERT INTO queue_items (room_id, video_id, position)
		VALUES ($1, $2, (SELECT COALESCE(MAX(position), 0) + 1
			FROM queue_items WHERE room_id = $2))
		RETURNING id, position
	`

	err := r.db.QueryRow(ctx, query, item.RoomID, item.VideoID).Scan(&item.ID, &item.Position)
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return item.ID, item.Position, nil
}

func (r *QueueItemRepository) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	const op = "repository.QueueItemRepository.Delete"

	if id == uuid.Nil {
		return -1, fmt.Errorf(
			"%s: %w: invalid queue item id", op, apperrors.ErrInvalidInput,
		)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return -1, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var roomID uuid.UUID
	var position int
	err = tx.QueryRow(
		ctx, `SELECT room_id, position FROM queue_items WHERE id = $1`, id,
	).Scan(&roomID, &position)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrQueueItemNotFound,
			)
		}
		return -1, fmt.Errorf("%s: get item: %w", op, err)
	}

	if _, err = tx.Exec(ctx, `DELETE FROM queue_items WHERE id = $1`, id); err != nil {
		return -1, fmt.Errorf("%s: delete failed: %w", op, err)
	}

	if _, err = tx.Exec(
		ctx,
		`UPDATE queue_items SET position = position - 1
		 WHERE room_id = $1 AND position > $2`,
		roomID, position,
	); err != nil {
		return -1, fmt.Errorf("%s: shift positions: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, fmt.Errorf("%s: commit: %w", op, err)
	}

	return position, nil
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
		SELECT id, room_id, video_id, is_active, position, created_at
		FROM queue_items
		WHERE id = $1
	`

	var item models.QueueItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.RoomID,
		&item.VideoID,
		&item.IsActive,
		&item.Position,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrQueueItemNotFound,
			)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &item, nil
}

func (r *QueueItemRepository) ListByRoom(
	ctx context.Context, roomID uuid.UUID,
	offset, limit int,
) ([]models.QueueItem, int, error) {
	const op = "repository.QueueItemRepository.ListByRoom"

	if roomID == uuid.Nil {
		return nil, 0, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
        SELECT id, room_id, video_id, is_active, position, created_at, 
               COUNT(*) OVER() as total_count
        FROM queue_items
        WHERE room_id = $1
        ORDER BY position ASC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	items := make([]models.QueueItem, 0, limit)
	total := 0

	for rows.Next() {
		var item models.QueueItem
		err := rows.Scan(
			&item.ID,
			&item.RoomID,
			&item.VideoID,
			&item.IsActive,
			&item.Position,
			&item.CreatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return items, total, nil
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
				op, id.String(), apperrors.ErrQueueItemNotFound,
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
