package order

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	ctx := context.Background()

	serviceOrder := getRepositoryCancelSuccess()
	order := updateRepositoryCancelSuccess()

	s.orderRepository.On("Get", ctx, "1").Return(serviceOrder, nil)
	s.orderRepository.On("Update", ctx, "1", order).Return(nil)
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

	order := getRepositoryCancelAlreadyPaid()

	s.orderRepository.On("Get", ctx, "1").Return(order, nil)
	err := s.service.Cancel(ctx, "1")

	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrAlreadyPaid)
}

func getRepositoryCancelAlreadyPaid() serviceModel.Order {
	return serviceModel.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          serviceModel.PAID,
	}
}

func getRepositoryCancelSuccess() serviceModel.Order {
	return serviceModel.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          serviceModel.CANCELLED,
	}
}

func updateRepositoryCancelSuccess() model.Order {
	return model.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          model.CANCELLED,
	}
}
