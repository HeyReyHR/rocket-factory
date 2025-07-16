package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/converter"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order By UUID " + params.OrderUUID + " not found",
		}, nil
	}
	return converter.ServiceOrderToRespOrder(order), nil
}
