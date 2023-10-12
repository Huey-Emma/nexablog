package permission

import (
	"context"
	"time"

	"nexablog/internal/repository/permission"
)

type Service interface {
	AddUserPermission(context.Context, int, ...string) error
}

type service struct {
	timeout time.Duration
	store   permission.Repo
}

func NewService(store permission.Repo) Service {
	return &service{
		3 * time.Second,
		store,
	}
}

func (s *service) AddUserPermission(ctx context.Context, userID int, codes ...string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.store.AddUserPermission(ctx, userID, codes...); err != nil {
		return err
	}

	return nil
}
