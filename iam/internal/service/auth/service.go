package auth

import (
	"time"

	"github.com/HeyReyHR/rocket-factory/iam/internal/repository"
	service2 "github.com/HeyReyHR/rocket-factory/iam/internal/service"
)

var _ service2.AuthService = (*service)(nil)

const (
	accessTokenSecret  = "access-secret-key-very-long-and-secure"
	refreshTokenSecret = "refresh-secret-key-very-long-and-secure"
	accessTokenTTL     = 15 * time.Minute
	refreshTokenTTL    = 24 * time.Hour
)

type service struct {
	repository     repository.AuthRepository
	userRepository repository.UserRepository
}

func NewService(repository repository.AuthRepository, userRepository repository.UserRepository) *service {
	return &service{
		repository:     repository,
		userRepository: userRepository,
	}
}
