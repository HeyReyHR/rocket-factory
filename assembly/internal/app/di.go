package app

import (
	"context"
	"fmt"

	"github.com/HeyReyHR/rocket-factory/assembly/internal/config"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/repository"
	assemblyRepository "github.com/HeyReyHR/rocket-factory/assembly/internal/repository/assembly"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/service"
	"github.com/HeyReyHR/rocket-factory/assembly/internal/service/assembly"
	orderConsumer "github.com/HeyReyHR/rocket-factory/assembly/internal/service/consumer/assembly_consumer"
	orderProducer "github.com/HeyReyHR/rocket-factory/assembly/internal/service/producer/assembly_producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/migrator"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	kafkaConverter "github.com/HeyReyHR/rocket-factory/assembly/internal/converter/kafka"
	wrappedKafka "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
)

type diContainer struct {
	assemblyService service.AssemblyService

	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer

	orderPaidDecoder       kafkaConverter.OrderPaidDecoder
	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.Producer

	assemblyRepository repository.AssemblyRepository

	postgresDBConn *pgx.Conn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderProducerService(ctx context.Context) service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderAssembledProducer(), d.AssemblyRepository(ctx))
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.AssemblyService(ctx))
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer() wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssembledProducer
}

func (d *diContainer) AssemblyService(ctx context.Context) service.AssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assembly.NewService(d.AssemblyRepository(ctx))
	}

	return d.assemblyService
}

func (d *diContainer) AssemblyRepository(ctx context.Context) repository.AssemblyRepository {
	if d.assemblyRepository == nil {
		d.assemblyRepository = assemblyRepository.NewRepository(d.PostgresDBConn(ctx))
	}

	return d.assemblyRepository
}

func (d *diContainer) PostgresDBConn(ctx context.Context) *pgx.Conn {
	if d.postgresDBConn == nil {
		dbConn, dbErr := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if dbErr != nil {
			panic(fmt.Sprintf("❌ failed to connect to Postgres: %s\n", dbErr.Error()))
		}

		dbErr = dbConn.Ping(ctx)
		if dbErr != nil {
			panic(fmt.Sprintf("❌ failed ping database: %s\n", dbErr.Error()))
		}

		migratorRunner := migrator.NewPgMigrator(stdlib.OpenDB(*dbConn.Config().Copy()), config.AppConfig().Postgres.MigrationsDir())
		dbErr = migratorRunner.Up()
		if dbErr != nil {
			logger.Error(ctx, "❌ failed to run migrations", zap.Error(dbErr))
			return nil
		}

		closer.AddNamed("Postgres conn", func(ctx context.Context) error {
			return dbConn.Close(ctx)
		})

		d.postgresDBConn = dbConn
	}

	return d.postgresDBConn
}
