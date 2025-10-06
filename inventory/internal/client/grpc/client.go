package grpc

import (
	"context"

	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc"
)

type IamClient interface {
	WhoAmI(ctx context.Context, in *authV1.WhoAmIRequest, opts ...grpc.CallOption) (*authV1.WhoAmIResponse, error)
	Login(ctx context.Context, in *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error)
}
