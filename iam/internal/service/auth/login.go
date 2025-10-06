package auth

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/iam/internal/repository/model"
	passwordManage "github.com/HeyReyHR/rocket-factory/iam/internal/utils/password"
	"github.com/google/uuid"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepository.GetWithLogin(ctx, login)
	if err != nil {
		return "", err
	}

	ok, err := passwordManage.VerifyArgon2id(password, user.PasswordHash)
	if !ok {
		return "", model.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	sessionUuid := uuid.NewString()
	now := time.Now().UTC()

	session := model.Session{
		Uuid:      sessionUuid,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(s.sessionTtl),
	}

	err = s.repository.Create(ctx, user, repoModel.Session(session), s.sessionTtl)
	if err != nil {
		return "", err
	}

	return sessionUuid, nil
}
