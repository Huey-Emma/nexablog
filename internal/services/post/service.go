package post

import (
	"context"
	"errors"
	"time"

	"nexablog/internal/models"
	"nexablog/internal/repository"
	"nexablog/internal/repository/post"
	"nexablog/internal/services"
)

type Service interface {
	CreatePost(context.Context, models.PostIn) (models.Post, error)
	FindAllPosts(context.Context) (models.Posts, error)
	FindPostByID(context.Context, int) (models.Post, error)
	DeletePostByID(context.Context, int) error
	UpdatePostByID(context.Context, models.PostIn, int, int) (models.Post, error)
	FindPostsByAuthor(context.Context, int) (models.Posts, error)
}

type service struct {
	timeout time.Duration
	store   post.Repo
}

func NewService(store post.Repo) Service {
	return &service{
		3 * time.Second,
		store,
	}
}

func (s *service) CreatePost(ctx context.Context, payload models.PostIn) (models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post, err := s.store.CreatePost(ctx, payload)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (s *service) FindAllPosts(ctx context.Context) (models.Posts, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	posts, err := s.store.FindAllPosts(ctx)
	if err != nil {
		return models.Posts{}, err
	}

	return posts, nil
}

func (s *service) FindPostByID(ctx context.Context, postID int) (models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post, err := s.store.FindPostByID(ctx, postID)

	if errors.Is(err, repository.ErrResourceNotFound) {
		return models.Post{}, services.ErrResourceNotFound
	}

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (s *service) DeletePostByID(ctx context.Context, postID int) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.store.DeletePostByID(ctx, postID)

	if errors.Is(err, repository.ErrResourceNotFound) {
		return services.ErrResourceNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePostByID(
	ctx context.Context,
	payload models.PostIn,
	postID, version int,
) (models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post, err := s.store.UpdatePostByID(ctx, payload, postID, version)

	if errors.Is(err, repository.ErrUpdateConflict) {
		return models.Post{}, services.ErrUpdateConflict
	}

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (s *service) FindPostsByAuthor(ctx context.Context, authorID int) (models.Posts, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	posts, err := s.store.FindPostsByAuthor(ctx, authorID)
	if err != nil {
		return models.Posts{}, err
	}

	return posts, nil
}
