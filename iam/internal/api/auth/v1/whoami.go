package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/iam/internal/converter"
	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
)

func (a *api) WhoAmI(ctx context.Context, r *authV1.WhoAmIRequest) (*authV1.WhoAmIResponse, error) {
	user, session, err := a.authService.WhoAmI(ctx, r.GetSessionUuid())
	if err != nil {
		return nil, err
	}
	return &authV1.WhoAmIResponse{
		Session: converter.ConvertSessionServiceToApi(session),
		User:    converter.ConvertUserServiceToApi(user),
	}, nil
}
