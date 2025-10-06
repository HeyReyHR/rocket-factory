package model

import (
	"errors"

	sharedErrors "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc/error"
)

var (
	ErrPartNotFound    = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrUnauthenticated = sharedErrors.NewUnauthenticatedError(errors.New("user not authenticated"))
)
