package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level           string `env:"LOGGER_LEVEL,required"`
	AsJson          bool   `env:"LOGGER_AS_JSON,required"`
	EnableOTLP      bool   `env:"ENABLE_OTLP,required"`
	OTLPServiceName string `env:"OTLP_SERVICE_NAME"`
	OTLPEnvironment string `env:"OTLP_SERVICE_ENVIRONMENT"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) EnableOTLP() bool {
	return cfg.raw.EnableOTLP
}

func (cfg *loggerConfig) OTLPServiceName() string {
	return cfg.raw.OTLPServiceName
}

func (cfg *loggerConfig) OTLPEnvironment() string {
	return cfg.raw.OTLPEnvironment
}
