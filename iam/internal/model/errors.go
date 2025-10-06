package model

import (
	"errors"

	sharedErrors "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc/error"
)

var (
	ErrUserNotFound           = sharedErrors.NewNotFoundError(errors.New("user not found"))
	ErrSessionNotFound        = sharedErrors.NewNotFoundError(errors.New("session not found"))
	ErrUserInvalidRegisterReq = sharedErrors.NewInvalidArgumentError(errors.New("invalid input"))
	ErrUserInvalidGetReq      = sharedErrors.NewInvalidArgumentError(errors.New("no uuid"))
	ErrInvalidCredentials     = sharedErrors.NewInvalidArgumentError(errors.New("invalid credentials"))
)
