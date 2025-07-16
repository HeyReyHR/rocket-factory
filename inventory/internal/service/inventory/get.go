package inventory

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

func (s *service) Part(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.inventoryRepository.Part(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
