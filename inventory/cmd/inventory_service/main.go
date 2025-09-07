package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/app"
	"github.com/HeyReyHR/rocket-factory/inventory/internal/config"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"go.uber.org/zap"
)

// const configPath = "deploy/compose/inventory/.env" // Ревьюверу: как мне избежать этого при коммите?

const configPath = "/home/heyrey/cool_projects/rocket-factory/deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Could not create app", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Error occurred while running app", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Error occurred while shutting down", zap.Error(err))
	}
}
