package e2e

import (
	"context"
	"os"
	"time"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/path"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Containers
	paymentAppName    = "payment-app"
	paymentDockerfile = "deploy/docker/payment/Dockerfile"

	// Env
	grpcPortKey = "PAYMENT_GRPC_PORT"

	// Env values
	loggerLevelValue = "debug"
	startupTimeout   = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Setting up test environment...")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "Cannot create network", zap.Error(err))
	}
	logger.Info(ctx, "✅ Network created")

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	projectRoot := path.GetProjectRoot()

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(paymentAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, paymentDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "Cannot run app container", zap.Error(err))
	}
	logger.Info(ctx, "✅ App container started")

	logger.Info(ctx, "🎉 Test environment is ready")
	return &TestEnvironment{
		Network: generatedNetwork,
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
