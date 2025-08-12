package e2e

import (
	"context"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUuid := uuid.NewString()

	return partUuid, nil

}

func (env *TestEnvironment) GetTestInventoryPart() *invV1.Part {
	return &invV1.Part{
		Uuid:          "engine-xxx",
		Name:          "Engine",
		Description:   "idk",
		Price:         150.01,
		Category:      invV1.Category_ENGINE,
		StockQuantity: 2,
		Manufacturer: &invV1.Manufacturer{
			Name:    "Bebra",
			Country: "France",
			Website: "govno.ru",
		},
		Tags: []string{"engine", "high-performance", "liquid"},
		Dimensions: &invV1.Dimensions{
			Length: 100.1,
			Width:  1.2,
			Height: 12.5,
			Weight: 500.1,
		},
		CreatedAt: timestamppb.New(time.Now().Add(-5 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now().Add(-1 * time.Hour)),
	}
}
