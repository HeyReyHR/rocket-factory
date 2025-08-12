package e2e

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Cleaning test environment...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ Test environment has been cleaned")
}

func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "Cannot stop app container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 App container stopped")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "Cannot delete network", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 network deleted")
		}
	}
}
