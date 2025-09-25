package user

import (
	"github.com/HeyReyHR/rocket-factory/iam/internal/repository"
	service2 "github.com/HeyReyHR/rocket-factory/iam/internal/service"
)

var _ service2.UserService = (*service)(nil)

type service struct {
	repository repository.UserRepository
}

func NewService(repository repository.UserRepository) *service {
	return &service{
		repository: repository,
	}
}
