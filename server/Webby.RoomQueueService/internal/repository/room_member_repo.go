package repository

import (
	"context"
	"fmt"
	"webby/room-queue-service/internal/apperrors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomMemberRepository struct {
	db *pgxpool.Pool
}

func NewRoomMemberRepository(db *pgxpool.Pool) *RoomMemberRepository {
	return &RoomMemberRepository{db: db}
}

func (r *RoomMemberRepository) Exists(
	ctx context.Context, roomId, userId uuid.UUID,
) (bool, error) {
	const op = "repository.RoomMemberRepository.Exists"

	if roomId == uuid.Nil {
		return false, fmt.Errorf(
			"%s: %w: invalid room id", op, apperrors.ErrInvalidInput,
		)
	}

	if userId == uuid.Nil {
		return false, fmt.Errorf(
			"%s: %w: invalid user id", op, apperrors.ErrInvalidInput,
		)
	}

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("EXISTS(SELECT 1 FROM room_members WHERE room_id = ? AND user_id = ?)")

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: build query: %w", op, err)
	}

	var exists bool
	err = r.db.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return exists, nil
}
