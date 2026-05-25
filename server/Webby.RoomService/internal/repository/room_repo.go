package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/room-service/internal/apperrors"
	"webby/room-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	const op = "repository.RoomRepository.Create"

	if room == nil {
		return uuid.Nil, fmt.Errorf("%s: %w: room cannot be nil", op, apperrors.ErrInvalidInput)
	}

	room.ID = uuid.New()

	query := `
		INSERT INTO rooms (id, host_id, category_id, name, thumbnail, is_private)
		VALUES ($1, $2, (SELECT id FROM categories WHERE name = $3), $4, $5, $6)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		room.ID,
		room.HostID,
		room.Category,
		room.Name,
		room.Thumbnail,
		room.IsPrivate,
	).Scan(&room.CreatedAt)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return room.ID, nil
}

func (r *RoomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.RoomRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM rooms WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	const op = "repository.RoomRepository.GetByID"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT r.id, r.host_id, c.name, r.name, r.thumbnail, r.is_private, r.created_at
		FROM rooms r
		JOIN categories c ON c.id = r.category_id
		WHERE r.id = $1
	`

	var room models.Room
	err := r.db.QueryRow(ctx, query, id).Scan(
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
			return nil, fmt.Errorf("%s: room %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &room, nil
}

func (r *RoomRepository) ListMy(ctx context.Context, userId uuid.UUID, page int, limit int, search string, categoryName *string) ([]models.Room, int64, error) {
	const op = "repository.RoomRepository.ListMy"

	if userId == uuid.Nil {
		return nil, 0, fmt.Errorf("%s: %w: invalid user id", op, apperrors.ErrInvalidInput)
	}

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

	whereClause := "WHERE r.host_id = $1"
	countArgs := []any{userId}
	paramN := 2

	if search != "" {
		whereClause += fmt.Sprintf(" AND r.name ILIKE $%d", paramN)
		countArgs = append(countArgs, "%"+search+"%")
		paramN++
	}

	if categoryName != nil {
		whereClause += fmt.Sprintf(" AND c.name = $%d", paramN)
		countArgs = append(countArgs, *categoryName)
		paramN++
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM rooms r JOIN categories c ON c.id = r.category_id %s`, whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.host_id, c.name, r.name, r.thumbnail, r.is_private, r.created_at
		FROM rooms r
		JOIN categories c ON c.id = r.category_id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramN, paramN+1)

	dataArgs := append(countArgs, limit, offset)

	rows, err := r.db.Query(ctx, query, dataArgs...)
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

func (r *RoomRepository) ListPublic(ctx context.Context, page int, limit int, search string, categoryName *string) ([]models.PublicRoom, int64, error) {
	const op = "repository.RoomRepository.ListPublic"

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

	whereClause := "WHERE r.is_private = false"
	countArgs := []any{}
	paramN := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND r.name ILIKE $%d", paramN)
		countArgs = append(countArgs, "%"+search+"%")
		paramN++
	}

	if categoryName != nil {
		whereClause += fmt.Sprintf(" AND c.name = $%d", paramN)
		countArgs = append(countArgs, *categoryName)
		paramN++
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM rooms r JOIN categories c ON c.id = r.category_id %s`, whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count query failed: %w", op, err)
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.host_id, u."Username", u."AvatarUrl", c.name, r.name, r.thumbnail, r.is_private, r.created_at
		FROM rooms r
		JOIN "Users" u ON u."UserId" = r.host_id
		JOIN categories c ON c.id = r.category_id
		%s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramN, paramN+1)

	dataArgs := append(countArgs, limit, offset)

	rows, err := r.db.Query(ctx, query, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: data query failed: %w", op, err)
	}
	defer rows.Close()

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

func (r *RoomRepository) Update(ctx context.Context, room *models.Room) (uuid.UUID, error) {
	const op = "repository.RoomRepository.Update"

	if room == nil || room.ID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		UPDATE rooms
		SET category_id = (SELECT id FROM categories WHERE name = $1), name = $2, thumbnail = $3, is_private = $4
		WHERE id = $5
	`

	tag, err := r.db.Exec(
		ctx,
		query,
		room.Category,
		room.Name,
		room.Thumbnail,
		room.IsPrivate,
		room.ID,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return uuid.Nil, fmt.Errorf("%s: room %s: %w", op, room.ID.String(), apperrors.ErrNotFound)
	}

	return room.ID, nil
}
