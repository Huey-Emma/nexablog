package post

import (
	"context"
	"database/sql"
	"errors"

	"nexablog/internal/models"
	"nexablog/internal/repository"
	"nexablog/internal/utils"
)

type Repo interface {
	CreatePost(context.Context, models.PostIn) (models.Post, error)
	FindAllPosts(context.Context) (models.Posts, error)
	FindPostByID(context.Context, int) (models.Post, error)
	DeletePostByID(context.Context, int) error
	UpdatePostByID(context.Context, models.PostIn, int, int) (models.Post, error)
	FindPostsByAuthor(context.Context, int) (models.Posts, error)
}

type repo struct {
	db utils.DBTX
}

func NewRepo(db utils.DBTX) Repo {
	return &repo{
		db,
	}
}

func (r *repo) CreatePost(ctx context.Context, payload models.PostIn) (models.Post, error) {
	q := `
  INSERT INTO posts 
  (title, body, author_id)
  VALUES ($1, $2, $3)
  RETURNING post_id, title, body, author_id, version, created_at;
  `

	row := r.db.QueryRowContext(ctx, q, payload.Title, payload.Body, payload.AuthorID)

	post := models.Post{}

	err := scanPost(row, &post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *repo) FindAllPosts(ctx context.Context) (models.Posts, error) {
	q := `
  SELECT post_id, title, body, author_id, version, created_at 
  FROM posts;
  `

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return make(models.Posts, 0), err
	}

	defer func() {
		_ = rows.Close()
	}()

	posts := make(models.Posts, 0)

	for rows.Next() {
		var post models.Post
		err := scanPost(rows, &post)
		if err != nil {
			return make(models.Posts, 0), err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return make(models.Posts, 0), err
	}

	return posts, nil
}

func (r *repo) FindPostByID(ctx context.Context, postID int) (models.Post, error) {
	q := `
  SELECT post_id, title, body, author_id, version, created_at
  FROM posts WHERE post_id = $1;
  `

	row := r.db.QueryRowContext(ctx, q, postID)

	post := models.Post{}

	err := scanPost(row, &post)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Post{}, repository.ErrResourceNotFound
	}

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *repo) DeletePostByID(ctx context.Context, postID int) error {
	q := `DELETE FROM posts WHERE post_id = $1;`

	result, err := r.db.ExecContext(ctx, q, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrResourceNotFound
	}

	return nil
}

func (r *repo) UpdatePostByID(
	ctx context.Context,
	payload models.PostIn,
	postID, version int,
) (models.Post, error) {
	q := `
  UPDATE posts SET title = $1, body = $2, version = version + 1
  WHERE post_id = $3 AND version = $4
  RETURNING post_id, title, body, author_id, version, created_at;
  `

	row := r.db.QueryRowContext(
		ctx,
		q,
		payload.Title,
		payload.Body,
		postID,
		version,
	)

	post := models.Post{}

	err := scanPost(row, &post)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Post{}, repository.ErrUpdateConflict
	}

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *repo) FindPostsByAuthor(ctx context.Context, authorID int) (models.Posts, error) {
	q := `
  SELECT post_id, title, body, author_id, version, created_at
  FROM posts WHERE author_id = $1;
  `

	rows, err := r.db.QueryContext(ctx, q, authorID)
	if err != nil {
		return make(models.Posts, 0), err
	}

	defer func() {
		_ = rows.Close()
	}()

	posts := make(models.Posts, 0)

	for rows.Next() {
		post := models.Post{}
		err := scanPost(rows, &post)
		if err != nil {
			return make(models.Posts, 0), err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return models.Posts{}, err
	}

	return posts, nil
}

func scanPost[R utils.Row](r R, p *models.Post) error {
	return r.Scan(
		&p.PostID,
		&p.Title,
		&p.Body,
		&p.AuthorID,
		&p.Version,
		&p.CreatedAt,
	)
}
