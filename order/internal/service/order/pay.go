package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/metrics"
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/tracing"
	gUuid "github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *service) Pay(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (string, error) {
	traceCtx, span := tracing.StartSpan(ctx, "order.pay",
		trace.WithAttributes(
			attribute.String("order.uuid", uuid),
		),
	)
	defer span.End()

	order, err := s.orderRepository.Get(traceCtx, uuid)
	if err != nil {
		span.RecordError(err)
		span.End()
		return "", model.ErrOrderNotFound
	}
	switch order.Status {
	case model.PAID:
		return "", model.ErrAlreadyPaid
	case model.CANCELLED:
		return "", model.ErrOrderCancelled
	case model.ASSEMBLED:
		return "", model.ErrOrderAlreadyAssembled
	default:
	}

	timeoutCtx, cancel := context.WithTimeout(traceCtx, model.RequestTimeout)
	defer cancel()

	transactionUuid, err := s.payment.PayOrder(timeoutCtx, uuid, paymentMethod)
	if err != nil {
		span.RecordError(err)
		span.End()
		return "", model.ErrPaymentNotProceeded
	}
	order.TransactionUuid = &transactionUuid

	order.PaymentMethod = &paymentMethod

	order.Status = model.PAID

	err = s.orderRepository.Update(traceCtx, uuid, converter.ServiceOrderToRepoOrder(order))
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	metrics.OrdersRevenueTotal.Add(traceCtx, order.TotalPrice)

	err = s.orderProducerService.ProduceOrderPaid(traceCtx, model.OrderPaidEvent{
		EventUuid:       gUuid.NewString(),
		OrderUuid:       uuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   string(paymentMethod),
		TransactionUuid: transactionUuid,
	})
	if err != nil {
		span.RecordError(err)
		span.End()
		return "", err
	}

	span.SetAttributes(
		attribute.String("payment.transactionUuid", transactionUuid),
		attribute.String("order.userUuid", order.UserUuid),
	)

	return transactionUuid, nil
}
