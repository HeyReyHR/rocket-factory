package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/converter"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, r *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	part, err := a.inventoryService.Part(ctx, r.GetUuid())
	if err != nil {
		return nil, err
	}
	return &invV1.GetPartResponse{
		Part: converter.PartModelToInfo(part),
	}, nil
}
