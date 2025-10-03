package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, info model.AdditionalInfo, password string) (string, error)
	User(ctx context.Context, uuid string) (model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	WhoAmI(ctx context.Context, sessionUuid string) (model.User, model.Session, error)
}
