package auth

import (
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository"
	service2 "github.com/HeyReyHR/rocket-factory/iam/internal/service"
)

var _ service2.AuthService = (*service)(nil)

type service struct {
	repository     repository.AuthRepository
	userRepository repository.UserRepository
	sessionTtl     time.Duration
}

func NewService(repository repository.AuthRepository, userRepository repository.UserRepository, sessionTtl time.Duration) *service {
	return &service{
		repository:     repository,
		userRepository: userRepository,
		sessionTtl:     sessionTtl,
	}
}
