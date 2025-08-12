package v1

import (
	"github.com/HeyReyHR/rocket-factory/payment/internal/service"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	service service.PaymentService
	payV1.UnimplementedPaymentServiceServer
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		service: paymentService,
	}
}
