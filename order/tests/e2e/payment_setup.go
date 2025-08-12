package e2e

import (
	"context"
	"os"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/testcontainers/network"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Containers
	paymentAppName    = "payment-app"
	paymentDockerfile = "deploy/docker/payment/Dockerfile"

	// Env
	paymentGrpcPortKey = "PAYMENT_GRPC_PORT"
)

func setupPaymentTestEnvironment(ctx context.Context, projectRoot string, generatedNetwork *network.Network) (*app.Container, error) {
	logger.Info(ctx, "🚀 Setting up payment...")

	grpcPort := getEnvWithLogging(ctx, paymentGrpcPortKey)

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
		logger.Fatal(ctx, "Cannot run payment container", zap.Error(err))
		return nil, err
	}
	logger.Info(ctx, "✅ Payment container started")

	return appContainer, nil
}
