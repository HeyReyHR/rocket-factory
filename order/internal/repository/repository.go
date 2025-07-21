package repository

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

type OrderRepository interface {
	Create(ctx context.Context, uuid, userUuid string, partUuids []string, totalPrice float64) (string, float64)
	Update(ctx context.Context, uuid string, order model.Order)
	Get(ctx context.Context, uuid string) (serviceModel.Order, error)
}
