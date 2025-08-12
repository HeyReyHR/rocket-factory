package inventory

import (
	"context"
	"fmt"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/model"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	var part model.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("findOne error: %e", err))

		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}
	logger.Debug(ctx, fmt.Sprint("part ", part))
	return part, nil
}
