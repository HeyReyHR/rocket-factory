package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func (r *repository) Create(_ context.Context, uuid, userUuid string, partUuids []string, totalPrice float64) (string, float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order := model.Order{
		Uuid:       uuid,
		UserUuid:   userUuid,
		PartUuids:  partUuids,
		TotalPrice: totalPrice,
		Status:     model.PENDING_PAYMENT,
	}
	r.data[uuid] = order
	return uuid, totalPrice
}
