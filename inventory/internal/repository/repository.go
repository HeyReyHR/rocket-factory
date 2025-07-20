package repository

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (serviceModel.Part, error)
	ListParts(ctx context.Context, filter serviceModel.Filter) ([]serviceModel.Part, error)
}
