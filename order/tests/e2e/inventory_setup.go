package e2e

import (
	"context"
	"os"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Containers
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Env
	inventoryGrpcPortKey = "INVENTORY_GRPC_PORT"
)

func setupInventoryTestEnvironment(ctx context.Context, projectRoot string, generatedNetwork *network.Network) (*mongo.Container, *app.Container, error) {

	logger.Info(ctx, "🚀 Setting up inventory...")

	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogging(ctx, inventoryGrpcPortKey)

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)

	if err != nil {
		logger.Fatal(ctx, "Cannot run MongoDB", zap.Error(err))
		return nil, nil, err
	}
	logger.Info(ctx, "✅ MongoDB container started")

	appEnv := map[string]string{
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)

	if err != nil {
		logger.Fatal(ctx, "Cannot run inventory container", zap.Error(err))
		return nil, nil, err
	}
	logger.Info(ctx, "✅ Inventory container started")

	return generatedMongo, appContainer, nil
}
