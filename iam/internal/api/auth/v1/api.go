package v1

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/service"
	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewApi(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
