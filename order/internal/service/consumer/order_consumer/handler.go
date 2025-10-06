package order_consumer

import (
	"context"
	"strconv"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/converter"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) ShipHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("build_time_sec", strconv.FormatInt(event.BuildTimeSec, 10)))

	order, err := s.repository.Get(ctx, event.OrderUuid)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}
	order.Status = model.ASSEMBLED
	err = s.repository.Update(ctx, order.Uuid, converter.ServiceOrderToRepoOrder(order))
	if err != nil {
		logger.Error(ctx, "Failed to update status", zap.Error(err))
		return err
	}

	return nil
}
