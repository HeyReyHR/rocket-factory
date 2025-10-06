package config

import (
	"os"

	"github.com/HeyReyHR/rocket-factory/inventory/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC ServiceConfig
	Mongo         MongoConfig
	IamClient     ServiceConfig
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

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	iamClientGRPCCfg, err := env.NewIamClientGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		Mongo:         mongoCfg,
		IamClient:     iamClientGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
