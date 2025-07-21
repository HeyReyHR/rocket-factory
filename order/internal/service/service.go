package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

type OrderService interface {
	Create(ctx context.Context, userUuid string, partUuids []string) (uuid string, totalPrice float64, err error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (transactionUuid string, err error)
	Cancel(ctx context.Context, uuid string) error
}

type OrderRepository interface {
	Create(ctx context.Context, uuid, userUuid string, partUuids []string, totalPrice float64) (string, float64, error)
	Update(ctx context.Context, uuid string, order repoModel.Order) error
	Get(ctx context.Context, uuid string) (model.Order, error)
}
