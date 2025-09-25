package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPServiceName() string
	OTLPEnvironment() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsDir() string
}
