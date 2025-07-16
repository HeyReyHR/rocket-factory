package v1

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.service.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order By UUID " + params.OrderUUID + " not found",
			}, nil
		case errors.Is(err, model.ErrAlreadyPaid):
			return &orderV1.ConfilctError{
				Code:    409,
				Message: "Order By UUID " + params.OrderUUID + " already paid, cannot be cancelled",
			}, nil
		}
	}
	return &orderV1.CancelOrderNoContent{}, nil
}
