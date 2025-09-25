package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	userV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, r *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userUuid, err := a.userService.Register(ctx, converter.ConvertUserInfoApiToService(r.GetInfo().GetInfo()), r.GetInfo().GetPassword())
	if err != nil {
		return nil, err
	}

	return &userV1.RegisterResponse{
		UserUuid: userUuid,
	}, nil
}
