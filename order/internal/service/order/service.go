package order

import (
	"github.com/HeyReyHR/rocket-factory/order/internal/client"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository"
	serviceInterface "github.com/HeyReyHR/rocket-factory/order/internal/service"
)

var _ serviceInterface.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	inventory       client.InventoryClient
	payment         client.PaymentClient
}

func NewService(inventory client.InventoryClient, payment client.PaymentClient, orderRepository repository.OrderRepository) *service {
	return &service{
		inventory:       inventory,
		payment:         payment,
		orderRepository: orderRepository,
	}
}
