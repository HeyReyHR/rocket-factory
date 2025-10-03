package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/metrics"
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
	gUuid "github.com/google/uuid"
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

	err = s.orderRepository.Update(ctx, uuid, converter.ServiceOrderToRepoOrder(order))
	if err != nil {
		return "", err
	}

	metrics.OrdersRevenueTotal.Add(ctx, order.TotalPrice)

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUuid:       gUuid.NewString(),
		OrderUuid:       uuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   string(paymentMethod),
		TransactionUuid: transactionUuid,
	})
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
