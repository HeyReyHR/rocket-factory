package e2e

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUuid := uuid.NewString()
	partDoc := bson.M{
		"uuid":           partUuid,
		"name":           gofakeit.Noun(),
		"description":    "High-performance rocket engine",
		"price":          gofakeit.Price(1000, 200000),
		"category":       "engine",
		"stock_quantity": 0,
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": "https://rocketcorp.com",
		},
		"tags": []string{"engine", "high-performance", "liquid"},
		"specs": bson.M{
			"tested":     true,
			"max_thrust": 25500,
			"fuel_type":  "liquid",
			"efficiency": 0.95,
		},
		"dimensions": bson.M{
			"length": gofakeit.Float64(),
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"weight": gofakeit.Float64(),
		},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "inventory-service"
	}

	_, err := env.Mongo.Client().Database(dbName).Collection(inventoryCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUuid, nil
}

func (env *TestEnvironment) InsertTestParts(ctx context.Context) ([]string, error) {
	partUuid1 := uuid.NewString()
	partUuid2 := uuid.NewString()

	partDoc1 := bson.M{
		"uuid":           partUuid1,
		"name":           gofakeit.Noun(),
		"description":    "High-performance rocket engine",
		"price":          gofakeit.Price(1000, 200000),
		"category":       "engine",
		"stock_quantity": 0,
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": "https://rocketcorp.com",
		},
		"tags": []string{"engine", "high-performance", "liquid"},
		"specs": bson.M{
			"tested":     true,
			"max_thrust": 25500,
			"fuel_type":  "liquid",
			"efficiency": 0.95,
		},
		"dimensions": bson.M{
			"length": gofakeit.Float64(),
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"weight": gofakeit.Float64(),
		},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	partDoc2 := bson.M{
		"uuid":           partUuid2,
		"name":           gofakeit.Noun(),
		"description":    "Advanced navigation system",
		"price":          gofakeit.Price(1000, 200000),
		"category":       "navigation",
		"stock_quantity": 0,
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": "https://rocketcorp.com",
		},
		"tags": []string{"navigation", "advanced", "system"},
		"specs": bson.M{
			"tested":     true,
			"accuracy":   "high",
			"range":      "global",
			"efficiency": 0.98,
		},
		"dimensions": bson.M{
			"length": gofakeit.Float64(),
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"weight": gofakeit.Float64(),
		},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "inventory-service"
	}

	documents := []interface{}{partDoc1, partDoc2}
	_, err := env.Mongo.Client().Database(dbName).Collection(inventoryCollectionName).InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}

	partUuids := []string{partUuid1, partUuid2}
	return partUuids, nil
}

func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
