package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUuid string, partUuids []string) (uuid string, totalPrice float64, err error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (transactionUuid string, err error)
	Cancel(ctx context.Context, uuid string) error
}

type ShipConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}
