package inventory

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestFilterSuccess() {
	ctx := context.Background()

	filter := model.Filter{
		Uuids: []string{"engine-001"},
	}

	part := []model.Part{
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

	s.inventoryRepository.On("Parts", ctx, filter).Return(part, nil)
	res, err := s.service.Parts(ctx, filter)
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestFilterEmptyResSuccess() {
	ctx := context.Background()

	filter := model.Filter{
		Uuids: []string{"engine-xxx"},
	}
	s.inventoryRepository.On("Parts", ctx, filter).Return([]model.Part{}, nil)
	res, err := s.service.Parts(ctx, filter)
	s.NoError(err)
	s.Equal([]model.Part{}, res)
}
