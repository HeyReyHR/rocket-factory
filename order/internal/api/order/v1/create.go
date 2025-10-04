package v1

import (
	"context"
	"errors"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	auth "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/grpc"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PostOrder(ctx context.Context, r *orderV1.CreateOrderRequest, params orderV1.PostOrderParams) (orderV1.PostOrderRes, error) {
	authorizedCtx := auth.ForwardSessionUUIDToGRPC(ctx)

	uuid, totalPrice, err := a.service.Create(authorizedCtx, r.UserUUID, r.PartUuids)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartOutOfStock):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Some parts are out of stock",
			}, nil
		case errors.Is(err, model.ErrListPartsFailed):
			return &orderV1.InternalServerError{
				Code:    500,
				Message: "Internal server error: " + err.Error(),
			}, err
		case errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Some parts haven't been found",
			}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		UUID:       uuid,
		TotalPrice: totalPrice,
	}, nil
}
