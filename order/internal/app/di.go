package app

import (
	"context"
	"fmt"

	orderV1API "github.com/HeyReyHR/rocket-factory/order/internal/api/order/v1"
	"github.com/HeyReyHR/rocket-factory/order/internal/client"
	iamClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/iam/v1"
	invClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/inventory/v1"
	payClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/payment/v1"
	"github.com/HeyReyHR/rocket-factory/order/internal/config"
	kafkaConverter "github.com/HeyReyHR/rocket-factory/order/internal/converter/kafka"
	"github.com/HeyReyHR/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository"
	orderRepository "github.com/HeyReyHR/rocket-factory/order/internal/repository/order"
	"github.com/HeyReyHR/rocket-factory/order/internal/service"
	shipConsumer "github.com/HeyReyHR/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/HeyReyHR/rocket-factory/order/internal/service/order"
	orderProducer "github.com/HeyReyHR/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/HeyReyHR/rocket-factory/platform/pkg/kafka/producer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/HeyReyHR/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/migrator"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/tracing"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	authV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/auth/v1"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	orderV1API orderV1.Handler

	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
	iamClient       client.IamClient

	orderProducerService service.OrderProducerService
	shipConsumerService  service.ShipConsumerService

	consumerGroup         sarama.ConsumerGroup
	shipAssembledConsumer wrappedKafka.Consumer

	orderService service.OrderService

	orderRepository repository.OrderRepository

	postgresDBConn *pgx.Conn

	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	syncProducer         sarama.SyncProducer
	orderProducer        wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewApi(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderRepository(ctx), d.OrderPaidProducerService())
	}

	return d.orderService
}

func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if d.paymentClient == nil {
		connPay, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("payment-service")))
		if err != nil {
			logger.Error(ctx, "❌ failed to connect to payment service", zap.Error(err))
		}

		payment := payV1.NewPaymentServiceClient(connPay)

		d.paymentClient = payClientV1.NewPaymentClient(payment)

		closer.AddNamed("Payment gRPC client", func(ctx context.Context) error {
			return connPay.Close()
		})
	}
	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if d.inventoryClient == nil {
		connInv, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error(ctx, "❌ failed to connect to inventory service", zap.Error(err))
			return nil
		}

		inventory := invV1.NewInventoryServiceClient(connInv)

		d.inventoryClient = invClientV1.NewInventoryClient(inventory)

		closer.AddNamed("Inventory gRPC client", func(ctx context.Context) error {
			return connInv.Close()
		})
	}
	return d.inventoryClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresDBConn(ctx))
	}

	return d.orderRepository
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

func (d *diContainer) OrderPaidProducerService() service.OrderProducerService {
	if d.orderProducer == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) ShipConsumerService(ctx context.Context) service.ShipConsumerService {
	if d.shipConsumerService == nil {
		d.shipConsumerService = shipConsumer.NewService(d.ShipConsumer(), d.ShipAssembledDecoder(), d.OrderRepository(ctx))
	}

	return d.shipConsumerService
}

func (d *diContainer) ShipConsumer() wrappedKafka.Consumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.shipAssembledConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.Config(),
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

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

func (d *diContainer) OrderProducer() wrappedKafka.Producer {
	if d.orderProducerService == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderProducer
}

func (d *diContainer) IamClient(ctx context.Context) client.IamClient {
	if d.iamClient == nil {
		connIam, err := grpc.NewClient(
			config.AppConfig().IamGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error(ctx, "❌ failed to connect to iam service", zap.Error(err))
			return nil
		}

		iam := authV1.NewAuthServiceClient(connIam)

		d.iamClient = iamClientV1.NewAuthClient(iam)
		closer.AddNamed("Iam gRPC client", func(ctx context.Context) error {
			return connIam.Close()
		})
	}
	return d.iamClient
}
