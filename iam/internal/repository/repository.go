package repository

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

type UserRepository interface {
	Create(ctx context.Context, info model.UserInfo, passwordHash string) (string, error)
	Get(ctx context.Context, uuid string) (model.User, error)
	GetWithLogin(ctx context.Context, login string) (model.User, error)
}

type AuthRepository interface {
	Create(ctx context.Context, user model.User) (string, error)
	Get(ctx context.Context, sessionUuid string) (model.User, model.Session, error)
}
