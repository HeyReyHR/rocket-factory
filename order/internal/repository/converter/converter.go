package converter

import (
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/order/internal/repository/model"
)

func RepoOrderToServiceOrder(order repoModel.Order) model.Order {
	return model.Order{
		Uuid:            order.Uuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   RepoPaymentMethodToService(order.PaymentMethod),
		Status:          RepoStatusToService(order.Status),
	}
}

func ServiceOrderToRepoOrder(order model.Order) repoModel.Order {
	return repoModel.Order{
		Uuid:            order.Uuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   ServicePaymentMethodToRepo(order.PaymentMethod),
		Status:          ServiceStatusToRepo(order.Status),
	}
}

func ServicePaymentMethodToRepo(paymentMethod *model.PaymentMethod) *repoModel.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}

	var result repoModel.PaymentMethod
	switch *paymentMethod {
	case model.CARD:
		result = repoModel.CARD
	case model.CREDIT_CARD:
		result = repoModel.CREDIT_CARD
	case model.SBP:
		result = repoModel.SBP
	case model.INVESTOR_MONEY:
		result = repoModel.INVESTOR_MONEY
	default:
		result = repoModel.UNKNOWN_METHOD
	}

	return &result
}

func RepoPaymentMethodToService(paymentMethod *repoModel.PaymentMethod) *model.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}

	var result model.PaymentMethod
	switch *paymentMethod {
	case repoModel.CARD:
		result = model.CARD
	case repoModel.CREDIT_CARD:
		result = model.CREDIT_CARD
	case repoModel.SBP:
		result = model.SBP
	case repoModel.INVESTOR_MONEY:
		result = model.INVESTOR_MONEY
	default:
		result = model.UNKNOWN_METHOD

	}

	return &result
}

func ServiceStatusToRepo(status model.Status) repoModel.Status {
	switch status {
	case model.PAID:
		return repoModel.PAID
	case model.CANCELLED:
		return repoModel.CANCELLED
	case model.PENDING_PAYMENT:
		return repoModel.PENDING_PAYMENT
	case model.ASSEMBLED:
		return repoModel.ASSEMBLED
	default:
		return repoModel.UNKNOWN_STATUS
	}
}

func RepoStatusToService(status repoModel.Status) model.Status {
	switch status {
	case repoModel.PAID:
		return model.PAID
	case repoModel.CANCELLED:
		return model.CANCELLED
	case repoModel.PENDING_PAYMENT:
		return model.PENDING_PAYMENT
	case repoModel.ASSEMBLED:
		return model.ASSEMBLED
	default:
		return model.UNKNOWN_STATUS
	}
}
