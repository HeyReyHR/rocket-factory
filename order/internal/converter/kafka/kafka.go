package kafka

import "github.com/HeyReyHR/rocket-factory/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
