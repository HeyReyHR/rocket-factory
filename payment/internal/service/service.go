package service

import (
	"context"

	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
)

type PaymentService interface {
	Pay(ctx context.Context, uuid string, paymentMethod payV1.PaymentMethod) (transactionUuid string, err error)
}
