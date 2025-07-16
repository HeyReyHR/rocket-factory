package order

import (
	"context"
	"log"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Create(ctx context.Context, userUuid string, partUuids []string) (string, float64, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, model.RequestTimeout)
	defer cancel()

	resp, err := s.inventory.ListParts(timeoutCtx, model.PartsFilter{
		Uuids: partUuids,
	})
	if err != nil {
		return "", 0, model.ErrListPartsFailed
	}

	if len(resp) != len(partUuids) {
		return "", 0, model.ErrPartsNotFound
	}

	for _, part := range resp {
		if part.StockQuantity == 0 {
			log.Printf("Part with uuid %s is out of stock", part.Uuid)
			return "", 0, model.ErrPartOutOfStock
		}
	}

	newUuid := uuid.NewString()

	var totalPrice float64
	for _, part := range resp {
		totalPrice += part.Price
	}

	s.orderRepository.Create(ctx, newUuid, userUuid, partUuids, totalPrice)

	return newUuid, totalPrice, nil
}
