package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"
	"webby/room-service/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *roomRepository {
	return &roomRepository{db}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	const op = "repository.roomRepository.Create"

	room.ID = uuid.New()

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("rooms").
		Columns("id", "host_id", "category_id", "name", "thumbnail", "is_private").
		Values(
			room.ID,
			room.HostID,
			sq.Expr("(SELECT id FROM categories WHERE name = ?)", room.Category),
			room.Name,
			room.Thumbnail,
			room.IsPrivate,
		).
		Suffix("RETURNING created_at").
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&room.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return room.ID, nil
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.roomRepository.Delete"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("rooms").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrRoomNotFound)
	}

	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	const op = "repository.roomRepository.GetByID"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("r.id", "r.host_id", "c.name", "r.name", "r.thumbnail", "r.is_private", "r.created_at").
		From("rooms r").
		Join("categories c ON c.id = r.category_id").
		Where(sq.Eq{"r.id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var room models.Room
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&room.ID,
		&room.HostID,
		&room.Category,
		&room.Name,
		&room.Thumbnail,
		&room.IsPrivate,
		&room.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrRoomNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &room, nil
}

func (r *roomRepository) ListMy(ctx context.Context, userID uuid.UUID, page, limit int, search, categoryName string) ([]models.Room, int64, error) {
	const op = "repository.roomRepository.ListMy"

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	countBuilder := builder.
		Select("COUNT(*)").
		From("rooms r").
		Join("categories c ON c.id = r.category_id").
		Where(sq.Eq{"r.host_id": userID})

	dataBuilder := builder.
		Select(
			"r.id",
			"r.host_id",
			"c.name",
			"r.name",
			"r.thumbnail",
			"r.is_private",
			"r.created_at",
		).
		From("rooms r").
		Join("categories c ON c.id = r.category_id").
		Where(sq.Eq{"r.host_id": userID}).
		OrderBy("r.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if search != "" {
		pattern := "%" + search + "%"
		countBuilder = countBuilder.Where(sq.Expr("r.name ILIKE ?", pattern))
		dataBuilder = dataBuilder.Where(sq.Expr("r.name ILIKE ?", pattern))
	}

	if categoryName != "" {
		countBuilder = countBuilder.Where(sq.Eq{"c.name": categoryName})
		dataBuilder = dataBuilder.Where(sq.Eq{"c.name": categoryName})
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var total int64
	err = r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

	rooms := []models.Room{}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.HostID,
			&room.Category,
			&room.Name,
			&room.Thumbnail,
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return rooms, total, nil
}

func (r *roomRepository) ListPublic(ctx context.Context, page, limit int, search, category string) ([]models.PublicRoom, int64, error) {
	const op = "repository.RoomRepository.ListPublic"
	log := logger.FromContext(ctx).With("op", op)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	countBuilder := builder.
		Select("COUNT(*)").
		From("rooms r").
		Join("categories c ON c.id = r.category_id").
		Where(sq.Eq{"r.is_private": false})

	dataBuilder := builder.
		Select(
			"r.id",
			"r.host_id",
			"u.\"Username\"",
			"u.\"AvatarUrl\"",
			"c.name",
			"r.name",
			"r.thumbnail",
			"r.is_private",
			"r.created_at",
		).
		From("rooms r").
		Join("\"Users\" u ON u.\"UserId\" = r.host_id").
		Join("categories c ON c.id = r.category_id").
		Where(sq.Eq{"r.is_private": false}).
		OrderBy("r.created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if search != "" {
		pattern := "%" + search + "%"
		countBuilder = countBuilder.Where(sq.Expr("r.name ILIKE ?", pattern))
		dataBuilder = dataBuilder.Where(sq.Expr("r.name ILIKE ?", pattern))
	}

	if category != "" {
		countBuilder = countBuilder.Where(sq.Eq{"c.name": category})
		dataBuilder = dataBuilder.Where(sq.Eq{"c.name": category})
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	var total int64
	err = r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	log.Debug("after query Count", "total", total)

	dataSQL, dataArgs, err := dataBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: build failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()
	log.Debug("after query select", "rows", rows)

	rooms := []models.PublicRoom{}
	for rows.Next() {
		var room models.PublicRoom
		err := rows.Scan(
			&room.ID,
			&room.HostID,
			&room.HostUsername,
			&room.HostAvatarUrl,
			&room.Category,
			&room.Name,
			&room.Thumbnail,
			&room.IsPrivate,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return rooms, total, nil
}

func (r *roomRepository) Update(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	const op = "repository.roomRepository.Update"

	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("rooms").
		Set("category_id", sq.Expr("(SELECT id FROM categories WHERE name = ?)", room.Category)).
		Set("name", room.Name).
		Set("thumbnail", room.Thumbnail).
		Set("is_private", room.IsPrivate).
		Where(sq.Eq{"id": room.ID}).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: build failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return uuid.Nil, fmt.Errorf("%s: %w", op, apperrors.ErrRoomNotFound)
	}

	return room.ID, nil
}

func (r *roomRepository) GetTotalRooms(ctx context.Context) (int, error) {
	const op = "repository.roomRepository.GetTotalRooms"

	sql, args, err := sq.Select("COUNT(id)").
		From("rooms").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: sql build failed: %w", op, err)
	}

	var activeRooms int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&activeRooms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: failed to execute: %w", op, err)
	}

	return activeRooms, nil
}
