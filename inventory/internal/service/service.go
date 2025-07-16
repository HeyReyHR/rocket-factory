package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

type InventoryService interface {
	Part(ctx context.Context, uuid string) (model.Part, error)
	Parts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}
