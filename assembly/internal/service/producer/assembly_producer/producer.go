package assembly_producer

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/repository"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type service struct {
	shipAssembledProducer kafka.Producer
	repository            repository.AssemblyRepository
}

func NewService(shipAssembledProducer kafka.Producer, repo repository.AssemblyRepository) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
		repository:            repo,
	}
}

func (p *service) ProcessAssembledEvents(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
			case <-ticker.C:
			}
			event, err := p.repository.Get(ctx)
			if err != nil {
				logger.Error(ctx, "failed to get new event", zap.Error(err))
				continue
			}
			if event.EventUuid == "" { // check if no pending events idk if it works
				continue
			}

			err = p.ProduceOrderAssembled(ctx, event)
			if err != nil {
				continue
			}

			if err := p.repository.Update(ctx, event.EventUuid); err != nil {
				logger.Error(ctx, "failed to set event done", zap.Error(err))
			}
		}
	}()
}

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {

	msg := &eventsV1.ShipAssembled{
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}
	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	err = p.shipAssembledProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembled", zap.Error(err))
		return err
	}

	return nil
}
