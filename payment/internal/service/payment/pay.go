package payment

import (
	"context"
	"log"

	"github.com/HeyReyHR/rocket-factory/payment/internal/model"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func (s *service) Pay(_ context.Context, orderUuid string, method payV1.PaymentMethod) (resUuid string, err error) {
	if orderUuid == "" {
		return "", model.ErrOrderUuidEmpty
	}
	if method == payV1.PaymentMethod_UNKNOWN {
		return "", model.ErrPaymentMethodUnknown
	}
	transactionUuid := uuid.NewString()
	log.Printf("Payment has been succeded, transaction_uuid: %s\n", transactionUuid)
	return transactionUuid, nil
}
