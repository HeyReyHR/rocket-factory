package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
)

func (s *service) Pay(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return "", model.ErrOrderNotFound
	}
	switch order.Status {
	case model.PAID:
		return "", model.ErrAlreadyPaid
	case model.CANCELLED:
		return "", model.ErrOrderCancelled
	default:
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, model.RequestTimeout)
	defer cancel()

	transactionUuid, err := s.payment.PayOrder(timeoutCtx, uuid, paymentMethod)
	if err != nil {
		return "", model.ErrPaymentNotProceeded
	}
	order.TransactionUuid = &transactionUuid

	order.PaymentMethod = &paymentMethod

	order.Status = model.PAID

	s.orderRepository.Update(ctx, uuid, converter.ServiceOrderToRepoOrder(order))

	return transactionUuid, nil
}
