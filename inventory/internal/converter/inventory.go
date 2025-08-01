package converter

import (
	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartModelToInfo(model model.Part) *invV1.Part {
	return &invV1.Part{
		Uuid:          model.Uuid,
		Name:          model.Name,
		Description:   model.Description,
		StockQuantity: model.StockQuantity,
		Dimensions:    DimensionsModelToInfo(model.Dimensions),
		CreatedAt:     timestamppb.New(model.CreatedAt),
		UpdatedAt:     timestamppb.New(model.UpdatedAt),
		Category:      CategoryModelToInfo(model.Category),
		Tags:          model.Tags,
		Manufacturer:  ManufacturerModelToInfo(model.Manufacturer),
		Metadata:      MetadataModelToInfo(model.Metadata),
		Price:         model.Price,
	}
}

func PartsModelToInfo(model []model.Part) []*invV1.Part {
	var partsInfo []*invV1.Part
	for _, part := range model {
		partsInfo = append(partsInfo, PartModelToInfo(part))
	}
	return partsInfo
}

func DimensionsModelToInfo(model model.Dimensions) *invV1.Dimensions {
	return &invV1.Dimensions{
		Length: model.Length,
		Width:  model.Width,
		Height: model.Height,
		Weight: model.Weight,
	}
}

func CategoryModelToInfo(category model.Category) invV1.Category {
	switch category {
	case model.ENGINE:
		return invV1.Category_ENGINE
	case model.FUEL:
		return invV1.Category_FUEL
	case model.PORTHOLE:
		return invV1.Category_PORTHOLE
	case model.WING:
		return invV1.Category_WING
	default:
		return invV1.Category_UNKNOWN
	}
}

func ManufacturerModelToInfo(model model.Manufacturer) *invV1.Manufacturer {
	return &invV1.Manufacturer{
		Name:    model.Name,
		Website: model.Website,
		Country: model.Country,
	}
}

func MetadataModelToInfo(model map[string]interface{}) map[string]*invV1.Value {
	if model == nil {
		return nil
	}
	infoMetadata := make(map[string]*invV1.Value)
	for key, value := range model {
		infoMetadata[key] = ValueModelToInfo(value)
	}
	return infoMetadata
}

func ValueModelToInfo(value interface{}) *invV1.Value {
	switch v := value.(type) {
	case string:
		return &invV1.Value{
			ValueType: &invV1.Value_StringValue{StringValue: v},
		}
	case int:
		return &invV1.Value{
			ValueType: &invV1.Value_Int64Value{Int64Value: int64(v)},
		}
	case int64:
		return &invV1.Value{
			ValueType: &invV1.Value_Int64Value{Int64Value: v},
		}
	case float32:
		return &invV1.Value{
			ValueType: &invV1.Value_DoubleValue{DoubleValue: float64(v)},
		}
	case float64:
		return &invV1.Value{
			ValueType: &invV1.Value_DoubleValue{DoubleValue: v},
		}
	case bool:
		return &invV1.Value{
			ValueType: &invV1.Value_BoolValue{BoolValue: v},
		}
	default:
		return &invV1.Value{}
	}
}

func CategoryInfoToModel(info invV1.Category) model.Category {
	switch info {
	case invV1.Category_ENGINE:
		return model.ENGINE
	case invV1.Category_FUEL:
		return model.FUEL
	case invV1.Category_PORTHOLE:
		return model.PORTHOLE
	case invV1.Category_WING:
		return model.WING
	default:
		return model.UNKNOWN
	}
}

func CategoriesInfoToModel(info []invV1.Category) []model.Category {
	var categories []model.Category
	for _, category := range info {
		categories = append(categories, CategoryInfoToModel(category))
	}
	return categories
}

func FilterInfoToModel(info *invV1.PartsFilter) model.Filter {
	return model.Filter{
		Uuids:                 info.Uuids,
		Names:                 info.Names,
		Tags:                  info.Tags,
		Categories:            CategoriesInfoToModel(info.Categories),
		ManufacturerCountries: info.ManufacturerCountries,
	}
}
