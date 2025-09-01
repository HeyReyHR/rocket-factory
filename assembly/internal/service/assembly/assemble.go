package assembly

import (
	"context"
	"time"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/assembly/internal/repository/model"
	"github.com/google/uuid"
)

func (s *service) Assemble(ctx context.Context, event model.OrderPaidEvent) error {
	delay := time.Duration(10) * time.Second
	time.Sleep(delay)

	eventUuid := uuid.NewString()

	err := s.assemblyRepository.Create(ctx, eventUuid, repoModel.EventType(model.OrderAssembledEventType), event.OrderUuid, event.UserUuid, int64(delay))
	if err != nil {
		return err
	}

	return nil
}
