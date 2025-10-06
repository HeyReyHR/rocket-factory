package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPServiceName() string
	OTLPEnvironment() string
}

type PaymentGRPCConfig interface {
	Address() string
}
