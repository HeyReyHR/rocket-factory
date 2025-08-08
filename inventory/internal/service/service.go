package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.Filter) ([]model.Part, error)
}
