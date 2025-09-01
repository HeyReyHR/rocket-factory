package order_consumer

import (
	"context"

	kafkaConverter "github.com/HeyReyHR/rocket-factory/notification/internal/converter/kafka"
	telegramService "github.com/HeyReyHR/rocket-factory/notification/internal/service"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        telegramService.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegram telegramService.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegram,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderAssembledConsumer service")
	err := s.orderAssembledConsumer.Consume(ctx, s.ShipHandler)

	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
