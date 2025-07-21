package order

import (
	"sync"

	repository2 "github.com/HeyReyHR/rocket-factory/order/internal/repository"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

var _ repository2.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]model.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]model.Order),
	}
}
