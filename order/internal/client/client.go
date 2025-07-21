package client

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)
}
