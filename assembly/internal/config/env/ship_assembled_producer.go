package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type shipAssembledProducerConfig struct {
	raw shipAssembledProducerEnvConfig
}

func NewShipAssembledProducerConfig() (*shipAssembledProducerConfig, error) {
	var raw shipAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledProducerConfig{raw: raw}, nil
}

func (cfg *shipAssembledProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *shipAssembledProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
