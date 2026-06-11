package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-queue-service/internal/apperrors"
	"webby/room-queue-service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queueItemRepository struct {
	db *pgxpool.Pool
}

func NewQueueItemRepository(db *pgxpool.Pool) *queueItemRepository {
	return &queueItemRepository{db}
}

func (r *queueItemRepository) Create(ctx context.Context, item *models.QueueItem) (uuid.UUID, int, error) {
	const op = "repository.queueItemRepository.Create"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("queue_items").
		Columns("room_id", "video_id", "position", "is_active").
		Values(
			item.RoomID,
			item.VideoID,
			sq.Expr("(SELECT COALESCE(MAX(position), 0) + 1 FROM queue_items WHERE room_id = ?)", item.RoomID),
			sq.Expr("(CASE WHEN (SELECT MAX(position) FROM queue_items WHERE room_id = ?) IS NULL THEN true ELSE false END)", item.RoomID),
		).
		Suffix("RETURNING id, position, is_active")

	sql, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("%s: build query: %w", op, err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&item.ID, &item.Position, &item.IsActive)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return uuid.Nil, -1, fmt.Errorf("%s: %w", op, apperrors.ErrConflict)
			}
		}
		return uuid.Nil, -1, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return item.ID, item.Position, nil
}

func (r *queueItemRepository) Delete(ctx context.Context, id uuid.UUID) (int, error) {
	const op = "repository.queueItemRepository.Delete"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return -1, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var roomID uuid.UUID
	var position int
	var isActive bool

	selectQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("room_id", "position", "is_active").
		From("queue_items").
		Where(sq.Eq{"id": id})

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s: build query: %w", op, err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&roomID, &position, &isActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, fmt.Errorf(
				"%s: queue item %s: %w",
				op, id.String(), apperrors.ErrQueueItemNotFound,
			)
		}
		return -1, fmt.Errorf("%s: get item: %w", op, err)
	}

	deleteQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete("queue_items").
		Where(sq.Eq{"id": id})

	sql, args, err = deleteQuery.ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s: build delete query: %w", op, err)
	}

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return -1, fmt.Errorf("%s: delete failed: %w", op, err)
	}

	shiftQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("queue_items").
		Set("position", sq.Expr("position - 1")).
		Where(sq.And{
			sq.Eq{"room_id": roomID},
			sq.Gt{"position": position},
		})

	sql, args, err = shiftQuery.ToSql()
	if err != nil {
		return -1, fmt.Errorf("%s: build shift query: %w", op, err)
	}

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return -1, fmt.Errorf("%s: shift positions: %w", op, err)
	}

	if isActive {
		activateQuery := sq.StatementBuilder.
			PlaceholderFormat(sq.Dollar).
			Update("queue_items").
			Set("is_active", true).
			Where(sq.And{
				sq.Eq{"room_id": roomID},
				sq.Eq{"position": position},
			})

		sql, args, err = activateQuery.ToSql()
		if err != nil {
			return -1, fmt.Errorf("%s: build activate query: %w", op, err)
		}

		cmdTag, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return -1, fmt.Errorf("%s: set active to next: %w", op, err)
		}

		if cmdTag.RowsAffected() == 0 {
			maxPosQuery := sq.StatementBuilder.
				PlaceholderFormat(sq.Dollar).
				Update("queue_items").
				Set("is_active", true).
				Where(sq.And{
					sq.Eq{"room_id": roomID},
					sq.Expr("position = (SELECT MAX(position) FROM queue_items WHERE room_id = ?)", roomID),
				})

			sql, args, err = maxPosQuery.ToSql()
			if err != nil {
				return -1, fmt.Errorf("%s: build max pos query: %w", op, err)
			}

			_, err = tx.Exec(ctx, sql, args...)
			if err != nil {
				return -1, fmt.Errorf("%s: set active to last: %w", op, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, fmt.Errorf("%s: commit: %w", op, err)
	}

	return position, nil
}

func (r *queueItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.QueueItem, error) {
	const op = "repository.queueItemRepository.GetByID"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "room_id", "video_id", "is_active", "position", "created_at").
		From("queue_items").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	var item models.QueueItem
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&item.ID,
		&item.RoomID,
		&item.VideoID,
		&item.IsActive,
		&item.Position,
		&item.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, apperrors.ErrQueueItemNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &item, nil
}

func (r *queueItemRepository) ListByRoom(ctx context.Context, roomID uuid.UUID, offset, limit int) ([]models.QueueItem, int, error) {
	const op = "repository.queueItemRepository.ListByRoom"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "room_id", "video_id", "is_active", "position", "created_at",
			"COUNT(*) OVER() as total_count").
		From("queue_items").
		Where(sq.Eq{"room_id": roomID}).
		OrderBy("position ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, -1, fmt.Errorf("%s: build query: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, -1, fmt.Errorf("%s: query failed: %w", op, err)
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
			return nil, -1, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, -1, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return items, total, nil
}

func (r *queueItemRepository) MoveToTop(ctx context.Context, id uuid.UUID) error {
	const op = "repository.queueItemRepository.MoveToTop"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var roomID uuid.UUID
	var currentPosition int

	selectQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("room_id", "position").
		From("queue_items").
		Where(sq.Eq{"id": id})

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&roomID, &currentPosition)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, apperrors.ErrQueueItemNotFound)
		}
		return fmt.Errorf("%s: get item: %w", op, err)
	}

	if currentPosition == 1 {
		return nil
	}

	shiftQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("queue_items").
		Set("position", sq.Expr("position + 1")).
		Where(sq.And{
			sq.Eq{"room_id": roomID},
			sq.Lt{"position": currentPosition},
		})

	sql, args, err = shiftQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: build shift query: %w", op, err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: shift positions: %w", op, err)
	}

	moveQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("queue_items").
		Set("position", 1).
		Where(sq.Eq{"id": id})

	sql, args, err = moveQuery.ToSql()
	if err != nil {
		return fmt.Errorf("%s: build move query: %w", op, err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: update position: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func (r *queueItemRepository) ActivateVideo(ctx context.Context, roomID, itemID uuid.UUID) (int, int, error) {
	const op = "repository.queueItemRepository.ActivateVideo"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	selectQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("position").
		From("queue_items").
		Where(sq.Eq{"id": itemID})

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return -1, -1, fmt.Errorf("%s: build query: %w", op, err)
	}

	var currentPosition int
	if err = tx.QueryRow(ctx, sql, args...).Scan(&currentPosition); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, -1, fmt.Errorf(
				"%s: queue item %s: %w",
				op, itemID.String(), apperrors.ErrQueueItemNotFound,
			)
		}
		return -1, -1, fmt.Errorf("%s: get item: %w", op, err)
	}

	updatePreviousQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("queue_items").
		Set("is_active", false).
		Where(sq.And{
			sq.Eq{"room_id": roomID},
			sq.Eq{"is_active": true},
			sq.NotEq{"id": itemID},
		}).
		Suffix("RETURNING position")

	sql, args, err = updatePreviousQuery.ToSql()
	if err != nil {
		return -1, -1, fmt.Errorf("%s: build update query: %w", op, err)
	}

	prevPosition := -1
	if err = tx.QueryRow(ctx, sql, args...).Scan(&prevPosition); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return -1, -1, fmt.Errorf("%s: update previous failed: %w", op, err)
		}
	}

	updateCurrentQuery := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Update("queue_items").
		Set("is_active", true).
		Where(sq.And{
			sq.Eq{"id": itemID},
			sq.Eq{"is_active": false},
		})
	sql, args, err = updateCurrentQuery.ToSql()
	if err != nil {
		return -1, -1, fmt.Errorf("%s: build update current query: %w", op, err)
	}

	commandTag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return -1, -1, fmt.Errorf("%s: set current active: %w", op, err)
	}
	if commandTag.RowsAffected() == 0 {
		return -1, -1, nil
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, -1, fmt.Errorf("%s: commit: %w", op, err)
	}

	return prevPosition, currentPosition, nil
}
