package model

import (
	"errors"

	sharedErrors "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc/error"
)

var (
	ErrOrderUuidEmpty       = sharedErrors.NewInvalidArgumentError(errors.New("order uuid is empty"))
	ErrPaymentMethodUnknown = sharedErrors.NewInvalidArgumentError(errors.New("payment method is unknown"))
)
