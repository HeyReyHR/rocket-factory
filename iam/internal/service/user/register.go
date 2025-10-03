package user

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	passwordManage "github.com/HeyReyHR/rocket-factory/iam/internal/utils/password"
)

func (s *service) Register(ctx context.Context, user model.AdditionalInfo, password string) (string, error) {
	if user.Email == "" || user.Login == "" {
		return "", model.ErrUserInvalidRegisterReq
	}

	passwordHash, err := passwordManage.HashArgon2id(password)
	if err != nil {
		return "", err
	}

	uuid, err := s.repository.Create(ctx, converter.ConvertUserInfoServiceToRepo(user), passwordHash)
	if err != nil {
		return "", err
	}
	return uuid, err
}
