package model

import "errors"

var (
	ErrOrderUuidEmpty       = errors.New("order uuid is empty")
	ErrPaymentMethodUnknown = errors.New("payment method is unknown")
)
