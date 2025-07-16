package v1

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/order/internal/client/converter"
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &invV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return converter.ProtoPartsToModel(parts.Parts), nil
}
