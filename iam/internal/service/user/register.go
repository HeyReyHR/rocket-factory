package user

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
)

func (s *service) Register(ctx context.Context, user model.UserInfo, password string) (string, error) {
	if user.Email == "" || user.Login == "" {
		return "", model.ErrUserInvalidRegisterReq
	}

	passwordHash, err := converter.HashArgon2id(password)
	if err != nil {
		return "", err
	}

	uuid, err := s.repository.Create(ctx, converter.ConvertUserInfoServiceToRepo(user), passwordHash)
	if err != nil {
		return "", err
	}
	return uuid, err
}
