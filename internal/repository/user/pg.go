package user

import (
	"context"
	"database/sql"
	"errors"

	"nexablog/internal/models"
	"nexablog/internal/repository"
	"nexablog/internal/utils"
)

type Repo interface {
	CreateUser(context.Context, models.UserIn) (models.User, error)
	FindUserByEmail(context.Context, string) (models.User, error)
	FindUserByToken(context.Context, string, models.Scope) (models.User, error)
}

type repo struct {
	db utils.DBTX
}

func NewRepo(db utils.DBTX) Repo {
	return &repo{
		db,
	}
}

func (r *repo) CreateUser(ctx context.Context, payload models.UserIn) (models.User, error) {
	q := `
  INSERT INTO users (username, email, password) 
  VALUES ($1, $2, $3)
  RETURNING user_id, username, email, password, version, created_at;
  `

	password, err := utils.GetPasswordHash(payload.Password)
	if err != nil {
		return models.User{}, err
	}

	row := r.db.QueryRowContext(
		ctx,
		q,
		payload.Username,
		payload.Email,
		password,
	)

	user := models.User{}

	err = scanUser(row, &user)

	if err != nil && repository.DuplicateKey(err) {
		return models.User{}, repository.ErrDuplicateKey
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *repo) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	q := `
  SELECT user_id, username, email, password, version, created_at
  FROM users WHERE email = $1;
  `

	row := r.db.QueryRowContext(ctx, q, email)

	var user models.User

	err := scanUser(row, &user)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, repository.ErrResourceNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func scanUser(row *sql.Row, u *models.User) error {
	return row.Scan(
		&u.UserID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Version,
		&u.CreatedAt,
	)
}

func (r *repo) FindUserByToken(
	ctx context.Context,
	token string,
	scope models.Scope,
) (models.User, error) {
	q := `
  SELECT 
  u.user_id, u.username, u.email, u.password, u.version, u.created_at
  FROM users u INNER JOIN tokens t USING(user_id)
  WHERE t.hash = $1 AND t.scope = $2 AND t.expires_at > now();
  `

	hash := utils.HashRandString(token)

	row := r.db.QueryRowContext(ctx, q, hash[:], scope)

	var user models.User

	err := scanUser(row, &user)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, repository.ErrResourceNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
