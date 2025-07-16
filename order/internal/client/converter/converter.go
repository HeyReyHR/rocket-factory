package converter

import (
	"github.com/HeyReyHR/rocket-factory/order/internal/model"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *invV1.PartsFilter {
	categories := make([]invV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		//nolint:gosec // category guaranteed to fit int32 range
		categories = append(categories, invV1.Category(category))
	}
	return &invV1.PartsFilter{
		Names:                 filter.Names,
		Categories:            categories,
		Uuids:                 filter.Uuids,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func ProtoMetadataToModel(metadata map[string]*invV1.Value) map[string]model.Value {
	result := make(map[string]model.Value, len(metadata))
	for key, value := range metadata {
		if value == nil || value.ValueType == nil {
			continue
		}
		switch v := value.ValueType.(type) {
		case *invV1.Value_StringValue:
			str := v.StringValue
			result[key] = model.Value{StringValue: &str}
		case *invV1.Value_Int64Value:
			i := v.Int64Value
			result[key] = model.Value{Int64Value: &i}
		case *invV1.Value_DoubleValue:
			f := v.DoubleValue
			result[key] = model.Value{DoubleValue: &f}
		case *invV1.Value_BoolValue:
			b := v.BoolValue
			result[key] = model.Value{BoolValue: &b}
		}
	}
	return result
}

func ProtoPartsToModel(parts []*invV1.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, model.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      model.Category(part.Category),
			Manufacturer: model.Manufacturer{
				Name:    part.Manufacturer.Name,
				Website: part.Manufacturer.Website,
				Country: part.Manufacturer.Country,
			},
			Tags:     part.Tags,
			Metadata: ProtoMetadataToModel(part.Metadata),
		})
	}
	return result
}
