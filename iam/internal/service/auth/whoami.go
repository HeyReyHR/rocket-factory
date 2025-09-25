package auth

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/converter"
)

func (s *service) WhoAmI(ctx context.Context, sessionUuid string) (model.User, model.Session, error) {
	user, session, err := s.repository.Get(ctx, sessionUuid)
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	return converter.ConvertUserRepoToService(user), converter.ConvertSessionRepoToService(session), nil
}
