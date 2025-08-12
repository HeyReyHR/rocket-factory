package config

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
}
