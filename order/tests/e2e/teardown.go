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

	if env.Postgres != nil {
		if err := env.Postgres.Connection().Close(ctx); err != nil {
			logger.Error(ctx, "Cannot stop Postgres container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Postgres container stopped")
		}
	}

	if env.Payment != nil {
		if err := env.Payment.Terminate(ctx); err != nil {
			logger.Error(ctx, "Cannot stop payment container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Payment container stopped")
		}
	}

	if env.Inventory != nil {
		if err := env.Inventory.Terminate(ctx); err != nil {
			logger.Error(ctx, "Cannot stop inventory container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Inventory container stopped")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "Cannot stop Mongo container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Mongo container stopped")
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
