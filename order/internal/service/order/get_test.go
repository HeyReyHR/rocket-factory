package order

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	ctx := context.Background()
	order := model.Order{
		Uuid:            "xxx",
		UserUuid:        "1",
		PartUuids:       []string{"engine-002", "fuel-001"},
		TotalPrice:      1000,
		TransactionUuid: nil,
		PaymentMethod:   nil,
		Status:          model.PENDING_PAYMENT,
	}
	s.orderRepository.On("Get", ctx, "xxx").Return(order, nil)
	res, err := s.service.Get(ctx, "xxx")
	s.NoError(err)
	s.Equal(order, res)
}

func (s *ServiceSuite) TestGetNotFoundError() {
	ctx := context.Background()
	s.orderRepository.On("Get", ctx, "x").Return(model.Order{}, model.ErrOrderNotFound)
	res, err := s.service.Get(ctx, "x")
	s.Error(err)
	s.Equal(err, model.ErrOrderNotFound)
	s.Empty(res)
}
