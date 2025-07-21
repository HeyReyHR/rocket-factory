package converter

import (
	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	repoModel "github.com/HeyReyHR/rocket-factory/inventory/internal/repository/model"
)

func CategoryModelToRepoModel(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func CategoryModelsToRepoModels(categories []model.Category) []repoModel.Category {
	var categoriesRes []repoModel.Category
	for _, category := range categories {
		categoriesRes = append(categoriesRes, CategoryModelToRepoModel(category))
	}
	return categoriesRes
}

func FilterModelToRepoModel(filter model.Filter) repoModel.Filter {
	return repoModel.Filter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoryModelsToRepoModels(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func RepoMetadataToMetadataModel(repoMetadata map[string]interface{}) map[string]interface{} {
	if repoMetadata == nil {
		return nil
	}

	modelMetadata := make(map[string]interface{})
	for key, value := range repoMetadata {
		modelMetadata[key] = value
	}
	return modelMetadata
}

func RepoModelToPartModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		Category:      model.Category(part.Category),
		StockQuantity: part.StockQuantity,
		Manufacturer:  RepoModelToManufacturerModel(part.Manufacturer),
		Tags:          part.Tags,
		CreatedAt:     part.CreatedAt,
		Dimensions:    RepoModelToDimensionsModel(part.Dimensions),
		UpdatedAt:     part.UpdatedAt,
		Metadata:      RepoMetadataToMetadataModel(part.Metadata),
	}
}

func RepoModelToDimensionsModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func RepoModelToManufacturerModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Website: manufacturer.Website,
		Country: manufacturer.Country,
	}
}

func RepoModelsToPartModels(parts []repoModel.Part) []model.Part {
	var partsRes []model.Part
	for _, part := range parts {
		partsRes = append(partsRes, RepoModelToPartModel(part))
	}
	return partsRes
}
