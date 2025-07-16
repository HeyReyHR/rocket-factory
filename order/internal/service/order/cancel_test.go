package order

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	ctx := context.Background()

	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          serviceModel.CANCELLED,
	}, nil)
	s.orderRepository.On("Update", ctx, "1", model.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          model.CANCELLED,
	})
	err := s.service.Cancel(ctx, "1")
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelNotFoundError() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{}, serviceModel.ErrOrderNotFound)
	err := s.service.Cancel(ctx, "1")
	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelAlreadyPaid() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          serviceModel.PAID,
	}, nil)
	err := s.service.Cancel(ctx, "1")
	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrAlreadyPaid)
}
