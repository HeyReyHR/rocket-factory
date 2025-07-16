package model

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrAlreadyPaid         = errors.New("order already paid")
	ErrOrderCancelled      = errors.New("order cancelled")
	ErrPaymentNotProceeded = errors.New("payment not proceeded")
	ErrListPartsFailed     = errors.New("list parts failed")
	ErrPartsNotFound       = errors.New("some parts not found")
	ErrPartOutOfStock      = errors.New("part out of stock")
)
