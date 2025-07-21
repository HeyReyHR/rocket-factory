package inventory

import (
	"context"
	"testing"
	"time"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/service/mocks"
)

func (s *ServiceSuite) TestsGetPart() {
	type inventoryServiceMockFunc func(t *testing.T) *mocks.InventoryService

	ctx := context.Background()

	tests := []struct {
		name                 string
		id                   string
		expectedRes          model.Part
		expectedErr          error
		inventoryServiceMock inventoryServiceMockFunc
	}{
		{
			name:        "GetPart success",
			id:          "engine-001",
			expectedRes: getRepositoryGetPartSuccess(),
			expectedErr: nil,
			inventoryServiceMock: func(t *testing.T) *mocks.InventoryService {
				mockService := mocks.NewInventoryService(t)
				mockService.On("GetPart", ctx, "engine-001").Return(getRepositoryGetPartSuccess(), nil)

				return mockService
			},
		},
		{
			name:        "GetPart not found err",
			id:          "engine-xxx",
			expectedRes: model.Part{},
			expectedErr: model.ErrPartNotFound,
			inventoryServiceMock: func(t *testing.T) *mocks.InventoryService {
				mockService := mocks.NewInventoryService(s.T())
				mockService.On("GetPart", ctx, "engine-xxx").Return(model.Part{}, model.ErrPartNotFound)

				return mockService
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			res, err := test.inventoryServiceMock(s.T()).GetPart(ctx, test.id)

			if test.expectedErr != nil {
				s.Error(err)
				s.Equal(test.expectedErr, err)
			} else {
				s.NoError(err)
			}

			s.Equal(test.expectedRes, res)
		})
	}
}

func getRepositoryGetPartSuccess() model.Part {
	return model.Part{
		Uuid:          "engine-001",
		Name:          "Rocket Engine V1",
		Description:   "High-performance rocket engine",
		Price:         15000.50,
		StockQuantity: 0,
		Category:      model.ENGINE,
		Manufacturer: model.Manufacturer{
			Name:    "RocketCorp",
			Country: "France",
			Website: "https://rocketcorp.com",
		},
		Tags: []string{"engine", "high-performance", "liquid"},
		Dimensions: model.Dimensions{
			Length: 2.5,
			Width:  1.0,
			Height: 1.0,
			Weight: 500.0,
		},
		CreatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
		UpdatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
	}
}
