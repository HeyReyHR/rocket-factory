package order

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateSuccess() {
	ctx := context.Background()
	partUuids := []string{"engine-002"}
	s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
		Uuids: partUuids,
	}).Return([]model.Part{{
		Uuid:          "engine-002",
		Name:          "Advanced Turbo Engine",
		Description:   "Next-generation turbo rocket engine",
		Price:         28500.00,
		StockQuantity: 5,
		Category:      model.ENGINE,
		Manufacturer: model.Manufacturer{
			Name:    "TurboTech",
			Country: "USA",
			Website: "https://turbotech.com",
		},
		Tags: []string{"engine", "turbo", "advanced", "high-thrust"},
		Dimensions: model.Dimensions{
			Length: 3.2,
			Width:  1.5,
			Height: 1.5,
			Weight: 750.0,
		},
		Metadata: map[string]model.Value{
			"max_thrust": {Int64Value: &[]int64{45000}[0]},
			"fuel_type":  {StringValue: &[]string{"hybrid"}[0]},
			"efficiency": {DoubleValue: &[]float64{0.98}[0]},
			"tested":     {BoolValue: &[]bool{true}[0]},
		},
		CreatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
		UpdatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
	}}, nil)
	s.orderRepository.On("Create", ctx, mock.AnythingOfType("string"), "heyrey", partUuids, 28500.00).Return("1", 28500.00)
	uuid, totalPrice, err := s.service.Create(ctx, "heyrey", partUuids)
	s.NoError(err)
	s.NotEmpty(uuid)
	s.Equal(28500.00, totalPrice)
}

func (s *ServiceSuite) TestCreateNotFoundParts() {
	ctx := context.Background()
	partUuids := []string{"engine-xxx"}
	s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
		Uuids: partUuids,
	}).Return([]model.Part{}, nil)
	uuid, totalPrice, err := s.service.Create(ctx, "heyrey", partUuids)
	s.Error(err)
	s.ErrorIs(err, model.ErrPartsNotFound)
	s.Empty(uuid)
	s.Equal(0.00, totalPrice)
}

func (s *ServiceSuite) TestCreatePartOutOfStock() {
	ctx := context.Background()
	partUuids := []string{"engine-002"}
	s.inventoryClient.On("ListParts", mock.Anything, model.PartsFilter{
		Uuids: partUuids,
	}).Return([]model.Part{{
		Uuid:          "engine-002",
		Name:          "Advanced Turbo Engine",
		Description:   "Next-generation turbo rocket engine",
		Price:         28500.00,
		StockQuantity: 0,
		Category:      model.ENGINE,
		Manufacturer: model.Manufacturer{
			Name:    "TurboTech",
			Country: "USA",
			Website: "https://turbotech.com",
		},
		Tags: []string{"engine", "turbo", "advanced", "high-thrust"},
		Dimensions: model.Dimensions{
			Length: 3.2,
			Width:  1.5,
			Height: 1.5,
			Weight: 750.0,
		},
		Metadata: map[string]model.Value{
			"max_thrust": {Int64Value: &[]int64{45000}[0]},
			"fuel_type":  {StringValue: &[]string{"hybrid"}[0]},
			"efficiency": {DoubleValue: &[]float64{0.98}[0]},
			"tested":     {BoolValue: &[]bool{true}[0]},
		},
		CreatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
		UpdatedAt: time.Date(2024, 11, 1, 10, 15, 30, 0, time.UTC),
	}}, nil)

	uuid, totalPrice, err := s.service.Create(ctx, "heyrey", partUuids)
	s.Error(err)
	s.ErrorIs(err, model.ErrPartOutOfStock)
	s.Empty(uuid)
	s.Equal(0.00, totalPrice)
}
