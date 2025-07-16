package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return model.ErrOrderNotFound
	}

	if order.Status == model.PAID {
		return model.ErrAlreadyPaid
	}

	order.Status = model.CANCELLED

	s.orderRepository.Update(ctx, uuid, converter.ServiceOrderToRepoOrder(order))

	return nil
}
