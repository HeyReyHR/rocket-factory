package repository

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	Part(ctx context.Context, uuid string) (serviceModel.Part, error)
	Parts(ctx context.Context, filter serviceModel.Filter) ([]serviceModel.Part, error)
}
