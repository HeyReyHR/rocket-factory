package app

import (
	"context"
	"fmt"

	orderV1API "github.com/HeyReyHR/rocket-factory/order/internal/api/order/v1"
	"github.com/HeyReyHR/rocket-factory/order/internal/client"
	invClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/inventory/v1"
	payClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/payment/v1"
	"github.com/HeyReyHR/rocket-factory/order/internal/config"
	"github.com/HeyReyHR/rocket-factory/order/internal/repository"
	orderRepository "github.com/HeyReyHR/rocket-factory/order/internal/repository/order"
	"github.com/HeyReyHR/rocket-factory/order/internal/service"
	orderService "github.com/HeyReyHR/rocket-factory/order/internal/service/order"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/closer"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/logger"
	"github.com/HeyReyHR/rocket-factory/platform/pkg/migrator"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
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

	orderService service.OrderService

	orderRepository repository.OrderRepository

	postgresDBConn *pgx.Conn
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
		d.orderService = orderService.NewService(d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderRepository(ctx))
	}

	return d.orderService
}

func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient { // Ревьюверу: контекст нужен или удалить нах?
	if d.paymentClient == nil {
		connPay, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
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

		migrationsDir := "migrations"

		migratorRunner := migrator.NewPgMigrator(stdlib.OpenDB(*dbConn.Config().Copy()), migrationsDir)

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
