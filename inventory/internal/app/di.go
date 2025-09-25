package app

import (
	"context"
	"fmt"

	invV1API "github.com/HeyReyHR/rocket-factory/inventory/internal/api/inventory/v1"
	client "github.com/HeyReyHR/rocket-factory/inventory/internal/client/grpc"
	iamClientV1 "github.com/HeyReyHR/rocket-factory/inventory/internal/client/grpc/iam/v1"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/config"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/repository"
	inventoryRepository "github.com/HeyReyHR/rocket-factory/inventory/internal/repository/inventory"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/service"
	inventoryService "github.com/HeyReyHR/rocket-factory/inventory/internal/service/inventory"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	inventoryV1API invV1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamClient client.IamClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) invV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = invV1API.NewApi(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		// nolint:contextcheck
		d.inventoryRepository = inventoryRepository.NewRepository(d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IamClient(ctx context.Context) client.IamClient {
	if d.iamClient == nil {
		connIam, err := grpc.NewClient(
			config.AppConfig().IamClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error(ctx, "❌ failed to connect to iam service", zap.Error(err))
			return nil
		}

		iam := authV1.NewAuthServiceClient(connIam)

		d.iamClient = iamClientV1.NewAuthClient(iam)
		closer.AddNamed("Iam gRPC client", func(ctx context.Context) error {
			return connIam.Close()
		})
	}
	return d.iamClient
}
