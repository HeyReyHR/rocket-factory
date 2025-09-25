package v1

import (
	"context"

	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, r *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionUuid, err := a.authService.Login(ctx, r.GetLogin(), r.GetPassword())
	if err != nil {
		return nil, err
	}

	return &authV1.LoginResponse{
		SessionUuid: sessionUuid,
	}, nil
}
