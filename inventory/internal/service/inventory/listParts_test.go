package inventory

import (
	"context"
	"testing"
	"time"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/service/mocks"
)

func (s *ServiceSuite) TestsList() {
	type inventoryServiceMockFunc func(t *testing.T) *mocks.InventoryService

	ctx := context.Background()

	parts, filter := getRepositoryListSuccess()

	emptyResFilter := getRepositoryListEmptyResSuccess()

	tests := []struct {
		name                 string
		filter               model.Filter
		expectedErr          error
		expectedRes          []model.Part
		inventoryServiceMock inventoryServiceMockFunc
	}{
		{
			name:        "List success",
			filter:      filter,
			expectedErr: nil,
			expectedRes: parts,
			inventoryServiceMock: func(t *testing.T) *mocks.InventoryService {
				mockService := mocks.NewInventoryService(t)
				mockService.On("ListParts", ctx, filter).Return(parts, nil)

				return mockService
			},
		},
		{
			name:        "List success empty",
			filter:      emptyResFilter,
			expectedErr: nil,
			expectedRes: []model.Part{},
			inventoryServiceMock: func(t *testing.T) *mocks.InventoryService {
				mockService := mocks.NewInventoryService(t)
				mockService.On("ListParts", ctx, filter).Return([]model.Part{}, nil)

				return mockService
			},
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			res, err := test.inventoryServiceMock(s.T()).ListParts(ctx, filter)

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

func getRepositoryListEmptyResSuccess() model.Filter {
	filter := model.Filter{
		Uuids: []string{"engine-xxx"},
	}
	return filter
}

func getRepositoryListSuccess() ([]model.Part, model.Filter) {
	filter := model.Filter{
		Uuids: []string{"engine-001"},
	}

	parts := []model.Part{
		{
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
			Metadata: map[string]model.Value{
				"max_thrust": {Int64Value: &[]int64{25000}[0]},
				"fuel_type":  {StringValue: &[]string{"liquid"}[0]},
				"efficiency": {DoubleValue: &[]float64{0.95}[0]},
				"tested":     {BoolValue: &[]bool{true}[0]},
			},
			CreatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
			UpdatedAt: time.Date(2024, 11, 3, 14, 27, 19, 0, time.UTC),
		},
	}
	return parts, filter
}
