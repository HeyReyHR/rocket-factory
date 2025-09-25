package user

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository/converter"
)

func (s *service) User(ctx context.Context, uuid string) (model.User, error) {
	if uuid == "" {
		return model.User{}, model.ErrUserInvalidGetReq
	}
	user, err := s.repository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return converter.ConvertUserRepoToService(user), err
}
