package payment

import (
	"context"

	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
)

func (s *ServiceSuite) TestPaySuccess() {
	ctx := context.Background()
	res, err := s.service.Pay(ctx, "123", payV1.PaymentMethod_CREDIT_CARD)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *ServiceSuite) TestPayEmptyOrderUuidFail() {
	ctx := context.Background()
	res, err := s.service.Pay(ctx, "", payV1.PaymentMethod_CREDIT_CARD)
	s.Error(err)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayUnknownPaymentMethodFail() {
	ctx := context.Background()
	res, err := s.service.Pay(ctx, "123", payV1.PaymentMethod_UNKNOWN)
	s.Error(err)
	s.Empty(res)
}
