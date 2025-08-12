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
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/postgres"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Containers
	orderAppName    = "order-app"
	orderDockerfile = "deploy/docker/order/Dockerfile"

	// Env
	httpPortKey = "HTTP_PORT"

	// Env values
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

type TestEnvironment struct {
	Network   *network.Network
	Postgres  *postgres.Container
	App       *app.Container
	Inventory *app.Container
	Mongo     *mongo.Container
	Payment   *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {

	logger.Info(ctx, "🚀 Setting up test environment...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "Cannot create network", zap.Error(err))
	}
	logger.Info(ctx, "✅ Network created")

	projectRoot := path.GetProjectRoot()

	generatedMongo, inventoryContainer, err := setupInventoryTestEnvironment(ctx, projectRoot, generatedNetwork)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
	}

	paymentContainer, err := setupPaymentTestEnvironment(ctx, projectRoot, generatedNetwork)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network:   generatedNetwork,
			Inventory: inventoryContainer,
			Mongo:     generatedMongo,
		})
	}

	postgresUsername := getEnvWithLogging(ctx, testcontainers.PostgresUsernameKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)

	httpPort := getEnvWithLogging(ctx, httpPortKey)

	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUsername, postgresPassword),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{
			Network:   generatedNetwork,
			Inventory: inventoryContainer,
			Mongo:     generatedMongo,
			Payment:   paymentContainer,
		})
		logger.Fatal(ctx, "Cannot run Postgres", zap.Error(err))
	}
	logger.Info(ctx, "✅ Postgres container started")

	appEnv := map[string]string{
		testcontainers.PostgresHostKey: generatedPostgres.Config().ContainerName,
	}

	waitStrategy := wait.ForListeningPort(nat.Port(httpPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(orderAppName),
		app.WithPort(httpPort),
		app.WithDockerfile(projectRoot, orderDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)

	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "Cannot run app container", zap.Error(err))
	}
	logger.Info(ctx, "✅ App container started")

	logger.Info(ctx, "🎉 Test environment is ready")
	return &TestEnvironment{
		Network:   generatedNetwork,
		Postgres:  generatedPostgres,
		App:       appContainer,
		Inventory: inventoryContainer,
		Mongo:     generatedMongo,
		Payment:   paymentContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Env is not set", zap.String("key", key))
	}

	return value
}
