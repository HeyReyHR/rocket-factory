package testcontainers

// DB constants
const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"

	PostgresContainerName = "postgres"
	PostgresPort          = "5432"

	PostgresImageNameKey = "POSTGRES_IMAGE_NAME"
	PostgresHostKey      = "POSTGRES_HOST"
	PostgresPortKey      = "POSTGRES_PORT"
	PostgresDatabaseKey  = "POSTGRES_DB"
	PostgresUsernameKey  = "POSTGRES_USER"
	// nolint:gosec
	PostgresPasswordKey = "POSTGRES_PASSWORD"

	InventoryServiceClientPortKey = "INVENTORY_CLIENT_GRPC_PORT"
	InventoryServiceClientHostKey = "INVENTORY_CLIENT_GRPC_HOST"

	PaymentServiceClientPortKey = "PAYMENT_CLIENT_GRPC_PORT"
	PaymentServiceClientHostKey = "PAYMENT_CLIENT_GRPC_HOST"
)
