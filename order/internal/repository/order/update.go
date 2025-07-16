package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func (r *repository) Update(_ context.Context, uuid string, order model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[uuid] = order
}
