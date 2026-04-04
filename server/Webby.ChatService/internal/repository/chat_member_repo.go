package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

	query := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return nil
}

func (r *ChatMemberRepository) Remove(ctx context.Context, chatId, userId uuid.UUID) error {
	const op = "repository.ChatMemberRepository.Remove"

	query := `DELETE FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	tag, err := r.db.Exec(ctx, query, chatId, userId)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: member not found: %w", op, fmt.Errorf("not found"))
	}

	return nil
}

func (r *ChatMemberRepository) Exists(ctx context.Context, chatId, userId uuid.UUID) (bool, error) {
	const op = "repository.ChatMemberRepository.Exists"

	query := `SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, chatId, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return exists, nil
}

func (r *ChatMemberRepository) ListByChat(ctx context.Context, chatId uuid.UUID) ([]uuid.UUID, error) {
	const op = "repository.ChatMemberRepository.ListByChat"

	query := `SELECT user_id FROM chat_members WHERE chat_id = $1`

	rows, err := r.db.Query(ctx, query, chatId)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var userIds []uuid.UUID
	for rows.Next() {
		var userId uuid.UUID
		if err := rows.Scan(&userId); err != nil {
			return nil, fmt.Errorf("%s: scan failed: %w", op, err)
		}
		userIds = append(userIds, userId)
	}

	return userIds, nil
}
