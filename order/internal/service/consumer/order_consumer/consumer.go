package order_consumer

import (
	"context"

	kafkaConverter "github.com/HeyReyHR/rocket-factory/order/internal/converter/kafka"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
	repository            repository.OrderRepository
}

func NewService(shipAssembledConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, repository repository.OrderRepository) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		repository:            repository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderAssembledConsumer service")
	err := s.shipAssembledConsumer.Consume(ctx, s.ShipHandler)

	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
