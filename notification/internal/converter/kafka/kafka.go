package kafka

import "github.com/HeyReyHR/rocket-factory/notification/internal/model"

type OrderPaidDecoder interface {
	DecodeOrderPaid(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledDecoder interface {
	DecodeOrderAssembled(data []byte) (model.OrderAssembledEvent, error)
}
