package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	serviceModel "github.com/HeyReyHR/rocket-factory/iam/internal/model"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.User, error) {
	var b []byte
	err := r.dbConn.QueryRow(ctx, getUserQ, uuid).Scan(&b)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, serviceModel.ErrUserNotFound
		}
		return model.User{}, err
	}

	var dto getUserDTO
	fmt.Println(string(b))
	if err = json.Unmarshal(b, &dto); err != nil {
		return model.User{}, err
	}
	out := model.User{
		Uuid:      dto.User.Uuid,
		CreatedAt: dto.User.CreatedAt,
		UpdatedAt: dto.User.UpdatedAt,
		UserInfo: model.UserInfo{
			Login:               dto.UserInfo.Login,
			Email:               dto.UserInfo.Email,
			NotificationMethods: dto.UserInfo.NotificationMethods,
		},
	}
	return out, nil
}

const getUserQ = `
SELECT json_build_object(
	'user', json_build_object(
		'uuid', u.uuid,
		'created_at', u.created_at,
		'updated_at', u.updated_at
	),
	'user_info', json_build_object(
		'login', ui.login,
		'email', ui.email,
		'notification_methods', COALESCE(
			json_agg(
				json_build_object(
					'provider_name', nm.provider_name,
					'target', nm.target
				)
			)FILTER (WHERE nm.provider_name IS NOT NULL),
			'[]'::json
		) 
	)
) AS result
FROM users u
JOIN user_infos ui
  ON ui.user_uuid = u.uuid
LEFT JOIN notification_methods nm
  ON nm.user_infos_uuid = ui.user_uuid
WHERE u.uuid = $1
GROUP BY u.uuid, u.created_at, u.updated_at, ui.login, ui.email;
	`

type getUserDTO struct {
	User struct {
		Uuid      string    `json:"uuid"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"user"`
	UserInfo struct {
		PasswordHash        string                     `json:"password_hash"`
		Login               string                     `json:"login"`
		Email               string                     `json:"email"`
		NotificationMethods []model.NotificationMethod `json:"notification_methods"`
	} `json:"user_info"`
}
