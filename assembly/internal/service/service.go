package service

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
)

type AssemblyService interface {
	Assemble(ctx context.Context, event model.OrderPaidEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProcessAssembledEvents(ctx context.Context, handlePeriod time.Duration)
	ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error
}
