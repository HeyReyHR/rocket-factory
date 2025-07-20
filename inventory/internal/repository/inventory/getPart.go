package inventory

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}
	return converter.RepoModelToPartModel(part), nil
}
