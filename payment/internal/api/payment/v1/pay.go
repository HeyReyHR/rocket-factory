package v1

import (
	"context"

	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, r *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error) {
	transactionUuid, err := a.service.Pay(ctx, r.OrderUuid, r.PaymentMethod)
	if err != nil {
		return nil, err
	}
	return &payV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
