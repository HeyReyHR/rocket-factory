package user

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/model"
)

func (s *service) User(ctx context.Context, uuid string) (model.User, error) {
	if uuid == "" {
		return model.User{}, model.ErrUserInvalidGetReq
	}
	user, err := s.repository.Get(ctx, uuid)
	if err != nil {
		
	}
}
