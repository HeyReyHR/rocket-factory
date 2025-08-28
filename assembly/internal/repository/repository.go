package repository

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/assembly/internal/repository/model"
)

type AssemblyRepository interface {
	Create(ctx context.Context, uuid string, eventType repoModel.EventType, orderUuid string, userUuid string, buildTimeSec int64) error
	Update(ctx context.Context, uuid string) error
	Delete(ctx context.Context, uuid string) error
	Get(ctx context.Context) (model.OrderAssembledEvent, error)
}
