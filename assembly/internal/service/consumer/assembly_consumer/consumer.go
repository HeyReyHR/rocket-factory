package assembly_consumer

import (
	"context"

	kafkaConverter "github.com/HeyReyHR/rocket-factory/assembly/internal/converter/kafka"
	assemblyService "github.com/HeyReyHR/rocket-factory/assembly/internal/service"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	assemblyService   assemblyService.AssemblyService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, assemblyService assemblyService.AssemblyService) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		assemblyService:   assemblyService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting assembly orderPaidConsumer service")
	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}
	return nil
}
