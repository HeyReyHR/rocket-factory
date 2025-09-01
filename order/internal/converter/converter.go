package converter

import (
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
)

func ReqPaymentMethodToServicePaymentMethod(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCREDITCARD:
		return model.CREDIT_CARD
	case orderV1.PaymentMethodCARD:
		return model.CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.INVESTOR_MONEY
	case orderV1.PaymentMethodSBP:
		return model.SBP
	default:
		return model.UNKNOWN_METHOD
	}
}

func ServiceOrderToRespOrder(order model.Order) *orderV1.OrderDto {
	var transactionUUID orderV1.OptString
	if order.TransactionUuid != nil {
		transactionUUID = orderV1.OptString{
			Value: *order.TransactionUuid,
			Set:   true,
		}
	} else {
		transactionUUID = orderV1.OptString{
			Set: false,
		}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.OptPaymentMethod{
			Value: ServicePaymentMethodToRespPaymentMethod(order.PaymentMethod),
			Set:   true,
		}
	} else {
		paymentMethod = orderV1.OptPaymentMethod{
			Set: false,
		}
	}

	return &orderV1.OrderDto{
		UUID:            order.Uuid,
		UserUUID:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          ServiceStatusToRespStatus(order.Status),
	}
}

func ServiceStatusToRespStatus(status model.Status) orderV1.OrderStatus {
	switch status {
	case model.PAID:
		return orderV1.OrderStatusPAID
	case model.CANCELLED:
		return orderV1.OrderStatusCANCELLED
	case model.ASSEMBLED:
		return orderV1.OrderStatusASSEMBLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func ServicePaymentMethodToRespPaymentMethod(method *model.PaymentMethod) orderV1.PaymentMethod {
	switch *method {
	case model.CREDIT_CARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.CARD:
		return orderV1.PaymentMethodCARD
	case model.INVESTOR_MONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	case model.SBP:
		return orderV1.PaymentMethodSBP
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
