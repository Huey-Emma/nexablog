package token

import (
	"context"
	"time"

	"nexablog/internal/models"
	"nexablog/internal/repository/token"
	"nexablog/internal/utils"
)

type Service interface {
	AddToken(context.Context, models.TokenIn) (models.TokenOut, error)
}

type service struct {
	timeout time.Duration
	store   token.Repo
}

func NewService(store token.Repo) Service {
	return &service{3 * time.Second, store}
}

func (s *service) AddToken(ctx context.Context, payload models.TokenIn) (models.TokenOut, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	plain, err := utils.GenRandString(16)
	if err != nil {
		return models.TokenOut{}, err
	}

	token := models.Token{
		Hash:      utils.HashRandString(plain),
		UserID:    payload.UserID,
		ExpiresAt: payload.ExpiresAt,
		Scope:     payload.Scope,
	}

	if err := s.store.CreateToken(ctx, token); err != nil {
		return models.TokenOut{}, err
	}

	out := models.TokenOut{
		Plain:     plain,
		ExpiresAt: payload.ExpiresAt,
	}

	return out, nil
}
