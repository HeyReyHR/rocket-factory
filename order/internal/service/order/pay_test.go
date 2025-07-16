package order

import (
	"context"

	serviceModel "github.com/HeyReyHR/rocket-factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestPaySuccess() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{
		Uuid:            "1",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          serviceModel.PENDING_PAYMENT,
	}, nil)
	s.paymentClient.On("PayOrder", mock.Anything, "1", serviceModel.CARD).Return("2", nil)
	s.orderRepository.On("Update", ctx, "1", mock.AnythingOfType("model.Order")).Return()

	uuid, err := s.service.Pay(ctx, "1", serviceModel.CARD)
	s.NoError(err)
	s.NotEmpty(uuid)
}

func (s *ServiceSuite) TestPayNotFoundError() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{}, serviceModel.ErrOrderNotFound)
	uuid, err := s.service.Pay(ctx, "1", serviceModel.CARD)
	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrOrderNotFound)
	s.Empty(uuid)
}

func (s *ServiceSuite) TestPayAlreadyPaidError() {
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
	uuid, err := s.service.Pay(ctx, "1", serviceModel.CARD)
	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrAlreadyPaid)
	s.Empty(uuid)
}

func (s *ServiceSuite) TestPayAlreadyCancelledError() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "1").Return(serviceModel.Order{
		Uuid:       "1",
		UserUuid:   "1",
		PartUuids:  []string{"engine-002", "fuel-001"},
		TotalPrice: 1000,
		Status:     serviceModel.CANCELLED,
	}, nil)
	uuid, err := s.service.Pay(ctx, "1", serviceModel.CARD)
	s.Error(err)
	s.ErrorIs(err, serviceModel.ErrOrderCancelled)
	s.Empty(uuid)
}
