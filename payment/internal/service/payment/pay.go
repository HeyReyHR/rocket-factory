package payment

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/payment/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/tracing"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *service) Pay(ctx context.Context, orderUuid string, method payV1.PaymentMethod) (resUuid string, err error) {
	ctx, span := tracing.StartSpan(ctx, "payment.pay",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUuid),
		),
	)
	defer span.End()
	
	if orderUuid == "" {
		span.RecordError(model.ErrOrderUuidEmpty)
		span.End()
		return "", model.ErrOrderUuidEmpty
	}
	if method == payV1.PaymentMethod_UNKNOWN {
		span.RecordError(model.ErrPaymentMethodUnknown)
		span.End()
		return "", model.ErrPaymentMethodUnknown
	}
	transactionUuid := uuid.NewString()
	logger.Info(ctx, "Payment has been succeeded", zap.String("transaction_uuid", transactionUuid))
	span.SetAttributes(
		attribute.String("payment.method", string(method)),
	)

	return transactionUuid, nil
}
