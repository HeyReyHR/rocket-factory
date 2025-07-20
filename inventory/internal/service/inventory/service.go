package inventory

import (
	def "github.com/HeyReyHR/rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository def.InventoryRepository
}

func NewService(inventoryRepository def.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
