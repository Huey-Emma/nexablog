package permission

import (
	"context"

	"github.com/lib/pq"

	"nexablog/internal/models"
	"nexablog/internal/utils"
)

type Repo interface {
	AddUserPermission(context.Context, int, ...string) error
	GetUserPermission(context.Context, int) (models.Permissions, error)
}

type repo struct {
	db utils.DBTX
}

func NewRepo(db utils.DBTX) Repo {
	return &repo{
		db,
	}
}

func (r *repo) AddUserPermission(ctx context.Context, userID int, codes ...string) error {
	q := `
  INSERT INTO users_permissions (user_id, permission_id) 
  (SELECT $1, permission_id FROM permissions WHERE code = ANY($2));
  `

	result, err := r.db.ExecContext(ctx, q, userID, pq.Array(codes))
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetUserPermission(ctx context.Context, userID int) (models.Permissions, error) {
	q := `
  SELECT p.permission_id, p.code FROM permissions p 
  INNER JOIN users_permissions up USING(permission_id)
  INNER JOIN users u USING(user_id)
  WHERE u.user_id = $1;
  `

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return make(models.Permissions, 0), err
	}

	defer func() {
		_ = rows.Close()
	}()

	permissions := make(models.Permissions, 0)

	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.PermissionID, &permission.Code)
		if err != nil {
			return make(models.Permissions, 0), err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return make(models.Permissions, 0), err
	}

	return permissions, nil
}
