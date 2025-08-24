package service

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
