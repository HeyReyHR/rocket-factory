package model

import (
	"errors"

	sharedErrors "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc/error"
)

var (
	ErrOrderNotFound         = sharedErrors.NewNotFoundError(errors.New("order not found"))
	ErrAlreadyPaid           = sharedErrors.NewInvalidArgumentError(errors.New("order already paid"))
	ErrOrderCancelled        = sharedErrors.NewInvalidArgumentError(errors.New("order cancelled"))
	ErrPaymentNotProceeded   = sharedErrors.NewInvalidArgumentError(errors.New("payment not proceeded"))
	ErrListPartsFailed       = sharedErrors.NewInvalidArgumentError(errors.New("list parts failed"))
	ErrPartsNotFound         = sharedErrors.NewNotFoundError(errors.New("some parts not found"))
	ErrPartOutOfStock        = sharedErrors.NewInvalidArgumentError(errors.New("part out of stock"))
	ErrOrderScanFailed       = sharedErrors.NewInvalidArgumentError(errors.New("order scan failed"))
	ErrOrderAlreadyAssembled = sharedErrors.NewInvalidArgumentError(errors.New("order already assembled"))
	ErrAuthFailed            = sharedErrors.NewInvalidArgumentError(errors.New("auth failed"))
)
