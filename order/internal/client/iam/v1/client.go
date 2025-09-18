package v1

import (
	"context"

	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/common/v1"
	"google.golang.org/grpc"
)

type client struct {
	generatedClient authV1.AuthServiceClient
}

func NewAuthClient(generatedClient authV1.AuthServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}

func (c client) Login(ctx context.Context, in *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error) {
	panic("implement me")
}

func (c client) WhoAmI(ctx context.Context, in *authV1.WhoAmIRequest, opts ...grpc.CallOption) (*authV1.WhoAmIResponse, error) {
	resp, err := c.generatedClient.WhoAmI(ctx, &authV1.WhoAmIRequest{
		SessionUuid: in.GetSessionUuid(),
	})
	if err != nil {
		return &authV1.WhoAmIResponse{
			Session: &commonV1.Session{},
			User:    &commonV1.User{},
		}, err
	}
	return &authV1.WhoAmIResponse{
		Session: resp.Session,
		User:    resp.User,
	}, nil
}
