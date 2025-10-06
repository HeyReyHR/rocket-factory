package config

import (
	"os"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	Postgres               PostgresConfig
	Metrics                MetricsConfig
	OrderAssembledProducer OrderAssembledProducerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
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

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := env.NewMetricsConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		Postgres:               postgresCfg,
		Metrics:                metricsCfg,
		OrderAssembledProducer: orderAssembledProducerCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
