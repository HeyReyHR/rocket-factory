package order

import (
	"context"
	"fmt"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) Create(ctx context.Context, userUuid string, partUuids []string) (string, float64, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, model.RequestTimeout)
	defer cancel()

	resp, err := s.inventory.ListParts(timeoutCtx, model.PartsFilter{
		Uuids: partUuids,
	})
	if err != nil {
		logger.Error(ctx, "list parts failed", zap.Error(err))
		return "", 0, model.ErrListPartsFailed
	}

	if len(resp) != len(partUuids) {
		return "", 0, model.ErrPartsNotFound
	}

	for _, part := range resp {
		fmt.Println(part)
		if part.StockQuantity == 0 {
			logger.Info(ctx, fmt.Sprintf("Part with uuid %s is out of stock", part.Uuid))
			return "", 0, model.ErrPartOutOfStock
		}
	}

	newUuid := uuid.NewString()

	var totalPrice float64
	for _, part := range resp {
		totalPrice += part.Price
	}

	_, _, err = s.orderRepository.Create(ctx, newUuid, userUuid, partUuids, totalPrice)
	if err != nil {
		return "", 0, err
	}
	return newUuid, totalPrice, nil
}
