package inventory

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestGetSuccess() {
	ctx := context.Background()

	part := model.Part{
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
	s.inventoryRepository.On("Part", ctx, "engine-001").Return(part, nil)
	res, err := s.service.Part(ctx, "engine-001")
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetNotFoundError() {
	ctx := context.Background()

	expectedErr := status.Errorf(codes.NotFound, "part with UUID not found")
	s.inventoryRepository.On("Part", ctx, "x").Return(model.Part{}, status.Errorf(codes.NotFound, "part with UUID not found"))
	res, err := s.service.Part(ctx, "x")
	s.Error(err)
	s.Equal(expectedErr, status.Errorf(codes.NotFound, "part with UUID not found"))
	s.Empty(res)
}
