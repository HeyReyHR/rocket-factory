package e2e

import (
	"context"
	"os"
	"time"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/path"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Containers
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Env
	grpcPortKey = "GRPC_PORT"

	// Env values
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Setting up test environment...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "Cannot create network", zap.Error(err))
	}
	logger.Info(ctx, "✅ Network created")

	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "Cannot run MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ MongoDB container started")

	projectRoot := path.GetProjectRoot()

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
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "Cannot run app container", zap.Error(err))
	}
	logger.Info(ctx, "✅ App container started")

	logger.Info(ctx, "🎉 Test environment is ready")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Env is not set", zap.String("key", key))
	}

	return value
}
