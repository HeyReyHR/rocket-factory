package v1

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/order/internal/converter"
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, r *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transactionUuid, err := a.service.Pay(ctx, params.OrderUUID, converter.ReqPaymentMethodToServicePaymentMethod(r.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderAlreadyAssembled):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Order with UUID" + params.OrderUUID + " already assembled",
			}, nil
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order with UUID " + params.OrderUUID + " not found",
			}, nil
		case errors.Is(err, model.ErrAlreadyPaid):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Order has already been paid",
			}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Cannot pay cancelled order",
			}, nil
		case errors.Is(err, model.ErrPaymentNotProceeded):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Could not proceed payment",
			}, nil
		}
	}

	return &orderV1.OrderPayResponse{
		TransactionUUID: transactionUuid,
	}, nil
}
