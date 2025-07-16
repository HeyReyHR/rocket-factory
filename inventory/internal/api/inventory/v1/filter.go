package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/converter"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, r *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.Parts(ctx, converter.FilterInfoToModel(r.Filter))
	if err != nil {
		return nil, err
	}
	return &invV1.ListPartsResponse{
		Parts: converter.PartsModelToInfo(parts),
	}, nil
}
