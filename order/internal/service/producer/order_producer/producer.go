package assembly_producer

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUuid,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
