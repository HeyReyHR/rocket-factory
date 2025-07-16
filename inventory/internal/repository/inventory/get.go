package inventory

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *repository) Part(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
	}
	return converter.RepoModelToPartModel(part), nil
}
