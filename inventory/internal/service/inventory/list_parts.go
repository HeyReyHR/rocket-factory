package inventory

import (
	"context"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}
	return parts, nil
}
