package config

import (
	"os"

	"github.com/HeyReyHR/rocket-factory/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	OrderHTTP     ServiceConfig
	PaymentGRPC   ServiceConfig
	InventoryGRPC ServiceConfig
	Postgres      PostgresConfig
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

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		OrderHTTP:     orderHTTPCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
		Postgres:      postgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
