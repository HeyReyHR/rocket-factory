package postgres

import (
	"context"

	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	Logger        Logger

	Host string
	Port string
}

func buildConfig(opts ...Option) *Config {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "postgres-container",
		ImageName:     "postgres:17",
		Username:      "test",
		Password:      "test",
		Database:      "db",
		Logger:        &logger.NoopLogger{},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
