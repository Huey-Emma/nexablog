package user

import (
	"context"
	"errors"
	"time"

	"nexablog/internal/models"
	"nexablog/internal/repository"
	"nexablog/internal/repository/user"
	"nexablog/internal/services"
)

type Service interface {
	CreateUser(context.Context, models.UserIn) (models.User, error)
	FindUserByEmail(context.Context, string) (models.User, error)
}

type service struct {
	timeout time.Duration
	store   user.Repo
}

func NewService(store user.Repo) Service {
	return &service{
		3 * time.Second,
		store,
	}
}

func (s *service) CreateUser(ctx context.Context, payload models.UserIn) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.store.CreateUser(ctx, payload)

	if errors.Is(err, repository.ErrDuplicateKey) {
		return models.User{}, services.ErrDuplicateKey
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *service) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.store.FindUserByEmail(ctx, email)

	if errors.Is(err, repository.ErrResourceNotFound) {
		return models.User{}, services.ErrResourceNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
