package assembly_producer

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembled{
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.shipAssembledProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
