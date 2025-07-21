package inventory

import (
	"context"
	"log"

	serviceModel "github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/converter"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, filter serviceModel.Filter) ([]serviceModel.Part, error) {
	filterQuery := bson.M{}

	if len(filter.Uuids) > 0 { // TODO FIX FILTER
		filterQuery["uuid"] = bson.M{
			"$in": filter.Uuids,
		}
	}
	if len(filter.Names) > 0 {
		filterQuery["name"] = bson.M{
			"$in": filter.Names,
		}
	}
	if len(filter.Categories) > 0 {
		filterQuery["category"] = bson.M{
			"$in": filter.Categories,
		}
	}
	if len(filter.ManufacturerCountries) > 0 {
		filterQuery["manufacturer.country"] = bson.M{
			"$in": filter.ManufacturerCountries,
		}
	}
	if len(filter.Tags) > 0 {
		filterQuery["tags"] = bson.M{
			"$in": filter.Tags,
		}
	}
	cursor, err := r.collection.Find(ctx, filterQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var parts []model.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	return converter.RepoModelsToPartModels(parts), nil
}
