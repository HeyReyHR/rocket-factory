package v1

import invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"

type client struct {
	generatedClient invV1.InventoryServiceClient
}

func NewInventoryClient(generatedClient invV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
