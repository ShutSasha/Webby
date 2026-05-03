package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	SaveToken(ctx context.Context, token string, userID uuid.UUID) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (svc *Service) GenerateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	const op = "services.GenerateToken"

	token := svc.generateToken()
	if err := svc.repository.SaveToken(ctx, token, userID); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (svc *Service) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	
	return base64.RawURLEncoding.EncodeToString(b)
}
