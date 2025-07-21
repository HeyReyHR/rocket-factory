package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, uuid string, method model.PaymentMethod) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, &payV1.PayOrderRequest{
		OrderUuid: uuid,
		//nolint:gosec // category guaranteed to fit int32 range
		PaymentMethod: payV1.PaymentMethod(method),
	})
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}
