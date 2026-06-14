package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

type tokenRepository interface {
	SaveToken(ctx context.Context, token string, userID uuid.UUID) error
	GetUserID(ctx context.Context, token string) (uuid.UUID, error)
}

type tokenService struct {
	repository tokenRepository
}

func NewTokenService(repository tokenRepository) *tokenService {
	return &tokenService{repository: repository}
}

func (svc *tokenService) GenerateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	const op = "services.GenerateToken"

	token := svc.generateToken()
	if err := svc.repository.SaveToken(ctx, token, userID); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (svc *tokenService) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	const op = "services.GetUserID"

	userID, err := svc.repository.GetUserID(ctx, token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (svc *tokenService) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)

	return base64.RawURLEncoding.EncodeToString(b)
}
