package inventory

import (
	"context"
	"slices"

	serviceModel "github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/converter"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) Parts(_ context.Context, filter serviceModel.Filter) ([]serviceModel.Part, error) {
	r.mu.RLock()
	var parts []model.Part
	for _, part := range r.data {
		parts = append(parts, model.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      part.Category,
			Manufacturer:  part.Manufacturer,
			Tags:          part.Tags,
			Metadata:      part.Metadata,
			Dimensions:    part.Dimensions,
			CreatedAt:     part.CreatedAt,
			UpdatedAt:     part.UpdatedAt,
		})
	}

	r.mu.RUnlock()

	filteredParts := r.filterParts(parts, converter.FilterModelToRepoModel(filter))
	return converter.RepoModelsToPartModels(filteredParts), nil
}

type filterT func(part *model.Part) bool

func (r *repository) filterParts(parts []model.Part, filter model.Filter) []model.Part {
	var result []model.Part

	filters := makeFilters(filter)

	for i := range parts {
		needAdd := true
		for _, filter := range filters {
			if !filter(&parts[i]) {
				needAdd = false
				break
			}
		}
		if needAdd {
			result = append(result, parts[i])
		}
	}
	return result
}

func makeFilters(filter model.Filter) []filterT {
	var filters []filterT
	if len(filter.Uuids) > 0 {
		filters = append(filters, func(part *model.Part) bool {
			return slices.Contains(filter.Uuids, part.Uuid)
		})
	}

	if len(filter.Names) > 0 {
		filters = append(filters, func(part *model.Part) bool {
			return slices.Contains(filter.Names, part.Name)
		})
	}

	if len(filter.Categories) > 0 {
		filters = append(filters, func(part *model.Part) bool {
			return slices.Contains(filter.Categories, part.Category)
		})
	}

	if len(filter.ManufacturerCountries) > 0 {
		filters = append(filters, func(part *model.Part) bool {
			return slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country)
		})
	}

	if len(filter.Tags) > 0 {
		filters = append(filters, func(part *model.Part) bool {
			for _, tag := range part.Tags {
				if slices.Contains(filter.Tags, tag) {
					return true
				}
			}
			return false
		})
	}

	return filters
}
