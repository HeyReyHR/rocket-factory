package order_consumer

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) ShipHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.DecodeOrderAssembled(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	if err = s.telegramService.SendOrderAssembledNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send notification", zap.Error(err))
	}
	return nil
}
