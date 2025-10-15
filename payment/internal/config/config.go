package config

import (
	"os"

	"github.com/HeyReyHR/rocket-factory/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
	Tracing     TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		Tracing:     tracingCfg,
		PaymentGRPC: paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
