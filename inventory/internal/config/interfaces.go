package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPServiceName() string
	OTLPEnvironment() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type ServiceConfig interface {
	Address() string
}
