package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type ServiceConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type ShipAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}
