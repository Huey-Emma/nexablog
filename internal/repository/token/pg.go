package token

import (
	"context"

	"nexablog/internal/models"
	"nexablog/internal/utils"
)

type Repo interface {
	CreateToken(context.Context, models.Token) error
}

type repo struct {
	db utils.DBTX
}

func NewRepo(db utils.DBTX) Repo {
	return &repo{
		db,
	}
}

func (r *repo) CreateToken(ctx context.Context, token models.Token) error {
	q := `
  INSERT INTO tokens (hash, user_id, expires_at, scope)
  VALUES ($1, $2, $3, $4)
  `

	result, err := r.db.ExecContext(
		ctx,
		q,
		token.Hash,
		token.UserID,
		token.ExpiresAt,
		token.Scope,
	)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}
