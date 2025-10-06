package v1

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/service"
	userV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewApi(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
