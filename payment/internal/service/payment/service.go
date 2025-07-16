package payment

import service2 "github.com/HeyReyHR/rocket-factory/payment/internal/service"

var _ service2.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
