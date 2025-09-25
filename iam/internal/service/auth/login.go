package auth

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login string, password string) (string, error) {

	user, err := s.userRepository.GetWithLogin(ctx, login)
	if err != nil {
		return "", err
	}

	ok, err := converter.VerifyArgon2id(password, user.PasswordHash)
	if !ok {
		return "", model.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	sessionUuid, err := s.repository.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return sessionUuid, nil
}
