package order_paid_consumer

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.DecodeOrderPaid(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}
	
	if err = s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send notification", zap.Error(err))
	}

	return nil
}
