package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderApiV1 "github.com/HeyReyHR/rocket-factory/order/internal/api/order/v1"
	inventoryClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/inventory/v1"
	paymentClientV1 "github.com/HeyReyHR/rocket-factory/order/internal/client/payment/v1"
	repoOrder "github.com/HeyReyHR/rocket-factory/order/internal/repository/order"
	serviceOrder "github.com/HeyReyHR/rocket-factory/order/internal/service/order"
	orderV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/openapi/order/v1"
	invV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/inventory/v1"
	payV1 "github.com/HeyReyHR/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	paymentServiceAddress   = "localhost:50052"
	inventoryServiceAddress = "localhost:50051"
	httpPort                = "8080"
	requestTimeout          = 5 * time.Second
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

func main() {
	connPay, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %s", err)
	}

	defer func() {
		if cerr := connPay.Close(); cerr != nil {
			log.Printf("failed to close connection with payment: %s", cerr)
		}
	}()

	payment := payV1.NewPaymentServiceClient(connPay)

	connInv, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service: %s", err)
	}

	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			log.Printf("failed to close connection with inventory: %s", cerr)
		}
	}()

	inventory := invV1.NewInventoryServiceClient(connInv)

	paymentClient := paymentClientV1.NewPaymentClient(payment)
	inventoryClient := inventoryClientV1.NewInventoryClient(inventory)

	repository := repoOrder.NewRepository()

	orderService := serviceOrder.NewService(inventoryClient, paymentClient, repository)

	orderApi := orderApiV1.NewApi(orderService)

	orderServer, err := orderV1.NewServer(orderApi)
	if err != nil {
		log.Printf("Error occured when creating OpenAPI server: %s", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 Starting server on port %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Error occurred when starting server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Error occurred when shutting down server: %s\n", err)
	}

	log.Println("✅ Server stopped")
}
