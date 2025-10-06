package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	userV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, r *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.User(ctx, r.GetUserUuid())
	if err != nil {
		return nil, err
	}

	return &userV1.GetUserResponse{
		User: converter.ConvertUserServiceToApi(user),
	}, nil
}
