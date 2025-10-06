package user

import (
	"context"
	"errors"

	serviceModel "github.com/HeyReyHR/rocket-factory/iam/internal/model"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/converter"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

const getUserWithLoginQ = `
SELECT 
    u.uuid,
    u.created_at,
    u.updated_at,
	ui.login,
	ui.email,
	ui.password_hash,
	 COALESCE(
        ARRAY_AGG(
            nm.provider_name || ':' || nm.target
        ) FILTER (WHERE nm.provider_name IS NOT NULL AND nm.target IS NOT NULL),
        ARRAY[]::text[]
    ) AS notification_methods
FROM users u
JOIN user_infos ui
  ON ui.user_uuid = u.uuid
LEFT JOIN notification_methods nm
  ON nm.user_infos_uuid = ui.user_uuid
WHERE ui.login = $1
GROUP BY u.uuid, u.created_at, u.updated_at, ui.login, ui.email, ui.password_hash;
	`

func (r *repository) GetWithLogin(ctx context.Context, login string) (model.User, error) {
	var dto model.UserDbDto
	var user model.User
	row := r.dbConn.QueryRow(ctx, getUserWithLoginQ, login)

	err := row.Scan(
		&dto.Uuid,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.Login,
		&dto.Email,
		&dto.PasswordHash,
		&dto.NotificationMethods,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, serviceModel.ErrUserNotFound
		}
		return model.User{}, err
	}
	if user, err = converter.ConvertUserDbToRepo(&dto); err != nil {
		return model.User{}, err
	}
	return user, nil
}
