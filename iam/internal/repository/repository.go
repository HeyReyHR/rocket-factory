package repository

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
)

type UserRepository interface {
	Create(ctx context.Context, info model.AdditionalInfo, passwordHash string) (string, error)
	Get(ctx context.Context, uuid string) (model.User, error)
	GetWithLogin(ctx context.Context, login string) (model.User, error)
}

type AuthRepository interface {
	Create(ctx context.Context, user model.User, session model.Session, sessionTtl time.Duration) error
	Get(ctx context.Context, sessionUuid string) (model.User, model.Session, error)
}
