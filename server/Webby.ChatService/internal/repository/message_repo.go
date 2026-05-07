package repository

import (
	"context"
	"errors"
	"fmt"
	"webby/chat-service/internal/apperrors"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *models.Message) (*models.Message, error) {
	const op = "repository.MessageRepository.Create"

	query := `
		INSERT INTO messages (sender_id, chat_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, msg.SenderID, msg.ChatID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return msg, nil
}

func (r *MessageRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	const op = "repository.MessageRepository.GetById"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid message id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, sender_id, chat_id, content, is_edited, edited_at, created_at
		FROM messages
		WHERE id = $1
	`

	var msg models.Message
	err := r.db.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.SenderID, &msg.ChatID,
		&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: message %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &msg, nil
}

func (r *MessageRepository) Update(ctx context.Context, id uuid.UUID, content string) (*models.Message, error) {
	const op = "repository.MessageRepository.Update"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid message id", op, apperrors.ErrInvalidInput)
	}

	query := `
		UPDATE messages
		SET content = $2, is_edited = true, edited_at = NOW()
		WHERE id = $1
		RETURNING id, sender_id, chat_id, content, is_edited, edited_at, created_at
	`

	var msg models.Message
	err := r.db.QueryRow(ctx, query, id, content).Scan(
		&msg.ID, &msg.SenderID, &msg.ChatID,
		&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: message %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return &msg, nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.MessageRepository.Delete"

	if id == uuid.Nil {
		return fmt.Errorf("%s: %w: invalid message id", op, apperrors.ErrInvalidInput)
	}

	query := `DELETE FROM messages WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: execution failed: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: message %s: %w", op, id.String(), apperrors.ErrNotFound)
	}

	return nil
}

func (r *MessageRepository) ListByChat(ctx context.Context, chatId uuid.UUID, page, limit int) ([]models.Message, int64, error) {
	const op = "repository.MessageRepository.ListByChat"

	if chatId == uuid.Nil {
		return nil, 0, fmt.Errorf("%s: %w: invalid chat id", op, apperrors.ErrInvalidInput)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM messages WHERE chat_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, chatId).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count failed: %w", op, err)
	}

	query := `
		SELECT id, sender_id, chat_id, content, is_edited, edited_at, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, chatId, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ChatID,
			&msg.Content, &msg.IsEdited, &msg.EditedAt, &msg.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%s: scan failed: %w", op, err)
		}
		messages = append(messages, msg)
	}

	return messages, total, nil
}
