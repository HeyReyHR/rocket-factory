package client

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)
}

type IamClient interface {
	WhoAmI(ctx context.Context, in *authV1.WhoAmIRequest, opts ...grpc.CallOption) (*authV1.WhoAmIResponse, error)
	Login(ctx context.Context, in *authV1.LoginRequest, opts ...grpc.CallOption) (*authV1.LoginResponse, error)
}
