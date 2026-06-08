package repository

import (
	"context"
	"fmt"
	"webby/chat-service/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatMemberRepository struct {
	db *pgxpool.Pool
}

func NewChatMemberRepository(db *pgxpool.Pool) *ChatMemberRepository {
	return &ChatMemberRepository{db: db}
}

func (r *ChatMemberRepository) Add(ctx context.Context, chatId, userId uuid.UUID) error {
	const op = "repository.ChatMemberRepository.Add"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Insert("chat_members").
		Columns("chat_id", "user_id").
		Values(chatId, userId).
		Suffix("ON CONFLICT (chat_id, user_id) DO NOTHING")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%s: query building failed: %w", op, err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return nil
}

func (r *ChatMemberRepository) Remove(ctx context.Context, chatId, userId uuid.UUID) error {
	const op = "repository.ChatMemberRepository.Remove"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Delete("chat_members").
		Where(sq.Eq{"chat_id": chatId, "user_id": userId})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("%s: query building failed: %w", op, err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: member not found: %w", op, fmt.Errorf("not found"))
	}

	return nil
}

func (r *ChatMemberRepository) Exists(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	const op = "repository.ChatMemberRepository.Exists"
	log := logger.FromContext(ctx).With("op", op)

	log.Debug("Parameters:", "chatID", chatID, "userID", userID)

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("1").
		From("chat_members").
		Where(sq.Eq{"chat_id": chatID, "user_id": userID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: query building failed: %w", op, err)
	}

	var exists int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: query failed: %w", op, err)
	}

	log.Debug("Passed exists check")

	return true, nil
}

func (r *ChatMemberRepository) ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error) {
	const op = "repository.ChatMemberRepository.ListByChat"

	query := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("user_id").
		From("chat_members").
		Where(sq.Eq{"chat_id": chatId})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: query building failed: %w", op, err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userId uuid.UUID
		if err := rows.Scan(&userId); err != nil {
			return nil, fmt.Errorf("%s: scan failed: %w", op, err)
		}
		userIDs = append(userIDs, userId)
	}

	return userIDs, nil
}
