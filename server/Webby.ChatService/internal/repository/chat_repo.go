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

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) (uuid.UUID, error) {
	const op = "repository.ChatRepository.Create"

	chat.Id = uuid.New()

	query := `
		INSERT INTO chats (id, room_id)
		VALUES ($1, $2)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query, chat.Id, chat.RoomId).Scan(&chat.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execution failed: %w", op, err)
	}

	return chat.Id, nil
}

func (r *ChatRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Chat, error) {
	const op = "repository.ChatRepository.GetById"

	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid chat id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, created_at
		FROM chats
		WHERE id = $1
	`

	var chat models.Chat
	err := r.db.QueryRow(ctx, query, id).Scan(&chat.Id, &chat.RoomId, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: chat %s: %w", op, id.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}

func (r *ChatRepository) GetByRoomId(ctx context.Context, roomId uuid.UUID) (*models.Chat, error) {
	const op = "repository.ChatRepository.GetByRoomId"

	if roomId == uuid.Nil {
		return nil, fmt.Errorf("%s: %w: invalid room id", op, apperrors.ErrInvalidInput)
	}

	query := `
		SELECT id, room_id, created_at
		FROM chats
		WHERE room_id = $1
	`

	var chat models.Chat
	err := r.db.QueryRow(ctx, query, roomId).Scan(&chat.Id, &chat.RoomId, &chat.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: room %s: %w", op, roomId.String(), apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return &chat, nil
}
